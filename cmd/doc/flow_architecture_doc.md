# Flow Architecture

## Обзор

Бот использует Finite State Machine (FSM) для управления многошаговыми диалогами с пользователями. Каждый пользователь имеет своё состояние, которое хранится в Redis.

## Компоненты системы

### 1. State (состояние)

```go
type State struct {
    Step ConversationStep  // Текущий шаг диалога
    Data StateData         // Данные, специфичные для flow
}
```

**Ключевая идея**: `StateData` - это union type (интерфейс с разными реализациями), каждый `Step` работает со своим типом данных.

### 2. StateData типы

| Тип данных | Используется в Steps                                                                                                                                                                                             | Назначение |
|-----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------|
| `IdleData` | `StepIdle`                                                                                                                                                                                                       | Пользователь не в диалоге |
| `SubscriptionCreationFlowData` | `StepAwaitingLabType`<br/>`StepAwaitingLabNumber`<br/>`StepAwaitingLabAuditorium`<br/>`StepAwaitingLabDomain`<br/>`SteAwaitingLabWeekday`<br/>`StepAwaitingLabLessons`<br/>`StepAwaitingSubCreationConfirmation` | Накапливает данные о создаваемой подписке |
| `SubscriptionListingFlowData` | `StepAwaitingListingSubsAction`                                                                                                                                                                                  | Хранит список подписок для навигации |

### 3. Router

Router - это middleware, который:
1. Извлекает `State` пользователя из Redis
2. Проверяет, является ли `update.Message.Text` командой, если да - сбрасывает состояние в `Idle` и выполняет команду
3. В противном случае находит handler для текущего `Step` в мапе `handlers[Step]`
4. Если handler найден - вызывает его с `StateData`
5. Если handler НЕ найден - сбрасывает состояние в `Idle`

**Важно**: Сброс состояния происходит автоматически, если для Step нет handler'а.

## Flows (потоки диалогов)

### Subscription Creation Flow

**Цель**: Создать новую подписку на уведомления о лабораторных работах.

**Steps в порядке выполнения**:

```
StepIdle
    ↓ (команда /sub)
StepAwaitingLabType
    ↓ (callback: performance/defence)
StepAwaitingLabNumber
    ↓ (текст: число)
    ├─→ StepAwaitingLabAuditorium (если performance)
    │       ↓ (текст: число)
    └─→ StepAwaitingLabDomain (если defence)
            ↓ (callback: mechanics/virtual/electricity)
            ↓
StepAwaitingLabWeekday
    ↓ (callback: день недели или skip)
StepAwaitingLabLessons
    ↓ (callback: номер пары или skip)
StepAwaitingSubCreationConfirmation
    ↓ (callback: create/cancel)
StepIdle
```

**StateData**: `SubscriptionCreationFlowData` - накапливает поля:

`UserID` → `LabType` → `LabNumber` → `LabAuditorium`/`LabDomain` → `Weekday` → `Lessons`

**Особенность**: Опциональные поля - pointer'ы.

### Subscription Listing Flow

**Цель**: Показать список подписок и дать возможность удалить.

**Steps**:

```
StepIdle
    ↓ (команда /list или /unsub)
StepAwaitingListingSubsAction
    ↓ (callback: move/delete)
    ├─→ остаёмся в StepAwaitingListingSubsAction (move - навигация)
    └─→ остаёмся в StepAwaitingListingSubsAction (delete)
```

**StateData**: `SubscriptionListingFlowData` - содержит:
- `UserSubs []ResponseSubscription` - полный список подписок пользователя

**Особенность**: Навигация по списку не меняет Step, только обновляет клавиатуру.

## Маппинг Step → StateData

Централизован в функции `dataTypeForStep()` в `fsm/state.go`:

```go
func dataTypeForStep(step ConversationStep) StateData {
    switch step {
    case StepIdle:
        return &IdleData{}
    case StepAwaitingLabType, StepAwaitingLabNumber, ...:
        return &SubscriptionCreationFlowData{}
    case StepAwaitingListingSubsAction:
        return &SubscriptionListingFlowData{}
    }
}
```

**Критически важно**: При добавлении нового Step обязательно обновить этот switch!

## Маппинг Step → Handler

Происходит при старте бота в `bot.go`:

```go
b.router.RegisterHandler(fsm.StepAwaitingLabType, b.handleLabType)
b.router.RegisterHandler(fsm.StepAwaitingLabNumber, b.handleLabNumber)
// ...
```

Хранится в `Router.handlers map[ConversationStep]HandlerFunc`.

## Переходы между Flow

**Правило**: При переходе из одного Flow в другой (или в Idle), старый State полностью заменяется новым.

**Примеры**:

1. Пользователь в Creation Flow → вводит `/list`:
   - Router замечает команду
   - Router сбрасывает State в `Idle`
   - Команда `/list` обрабатывается обычным handler'ом
   - Handler создаёт новый State с `SubscriptionListingFlowData`

2. Пользователь подтверждает создание подписки:
   - Handler явно вызывает `TryTransition(ctx, userID, StepIdle, &IdleData{})`
   - Старый `SubscriptionCreationFlowData` теряется

**Важно**: Переходы между Flow всегда явные - либо через команду (автоматический сброс), либо через `TryTransition` в handler'е.

## Десериализация State из Redis

**Проблема**: JSON хранит только данные, но не информацию о типе для интерфейса `StateData`.

**Решение**: Two-phase deserialization в `FSM.GetState()`:

1. Десериализуем в wrapper с `json.RawMessage`:
```go
var wrapper struct {
    Step ConversationStep
    Data json.RawMessage  // ещё не десериализовано
}
```

2. По `Step` определяем нужный тип через `dataTypeForStep()`

3. Десериализуем `Data` в конкретный тип:
```go
stateData := dataTypeForStep(wrapper.Step)
json.Unmarshal(wrapper.Data, stateData)
```

**Следствие**: `dataTypeForStep()` - критическая функция для корректной работы FSM.

## Обработка ошибок type assertion

В каждом handler'е происходит type assertion:
```go
flowData, ok := data.(*SubscriptionCreationFlowData)
if !ok {
    // Логируем critical error
    // Сбрасываем состояние в Idle
    // Показываем generic error пользователю
    return
}
```

**Когда это может случиться**:
- Баг в `dataTypeForStep()` (неправильный маппинг)
- Повреждённые данные в Redis
- Несовпадение версий кода (старые данные, новая структура)

**Стратегия**: Fail-safe - сбрасываем состояние и даём пользователю начать заново.

## Диаграмма потока данных

```
User Update
    ↓
[typingMiddleware] - показывает "typing..."
    ↓
[Router.Middleware]
    ├─→ Извлечь userID
    ├─→ Проверить не команда ли update.Message.Text
    │       └─→ Да: Сброс в Idle, выполнение команды
    ├─→ FSM.GetState(userID)
    │       ├─→ Redis.Get("fsm:{userID}:state")
    │       ├─→ json.Unmarshal (two-phase)
    │       └─→ State{Step, Data}
    ├─→ Найти handler для Step в map
    ├─→ Если найден:
    │       └─→ handler(ctx, bot, update, Data)
    │               ├─→ Type assertion Data → ConcreteType
    │               ├─→ Обработка логики
    │               ├─→ TryTransition(newStep, newData)
    │               └─→ SendMessage(response)
    └─→ Если НЕ найден:
            ├─→ FSM.ResetState(userID) → State{Idle, IdleData{}}
            └─→ next(ctx, bot, update) → обработка команды
```

## Добавление нового Flow - Checklist

1. ✅ Создать новый тип `StateData` в `fsm/state.go`
2. ✅ Добавить константы `Step` в `fsm/state.go`
3. ✅ Обновить `dataTypeForStep()` для новых Step'ов
4. ✅ Создать файл handler'ов (например `my_flow.go`)
5. ✅ Зарегистрировать handler'ы в `bot.Start()`
6. ✅ Обновить эту документацию!

## Частые вопросы

**Q: Почему StateData - интерфейс, а не generic?**
**A**: Map в Go не может хранить generic handlers с разными типами параметров. Интерфейс позволяет хранить все handler'ы в одной мапе.

**Q: Можно ли сделать несколько Flow одновременно?**
**A**: Нет, у пользователя всегда только один активный State. Переход в новый Flow сбрасывает предыдущий.

**Q: Что если пользователь не закончил Flow и закрыл бота?**
**A**: State хранится в Redis с TTL 24 часа. При возвращении пользователь продолжит с того же места (если не истёк TTL).

**Q: Зачем нужен StepIdle?**
**A**: Это default состояние. Когда пользователь не в диалоге, его State = Idle. Это позволяет отличить "нет состояния" от "ошибка получения состояния".