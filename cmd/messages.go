package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/metrics"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

func startMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>👋 Привет!</b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString("<b>Я бот для записи на лабораторные работы</b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString("<b>Буду следить за появлением доступных записей и сразу уведомлять тебя, когда появится нужная</b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString("<b>Используй /help для просмотра доступных команд</b>")
	return sb.String()
}

func helpMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>📖 Справка</b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString("<b>📝 Подписка:</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>/sub &lt;номер лабы&gt; &lt;номер аудитории&gt;</b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString("<b>⚙️ Управление:</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>/unsub &lt;номер подписки в списке&gt; - отписаться</b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString("<b>/list - посмотреть подписки</b>")
	return sb.String()
}

func subInvalidArgumentsMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Некорректные аргументы.</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>Использование: /sub &lt;номер лабы&gt; &lt;номер аудитории&gt;</b>")
	return sb.String()
}

func subLabNumberValidationErrorMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Номер лабы должен быть числом в диапазоне от 1 до 999</b>")
	return sb.String()
}

func subAuditoriumNumberValidationErrorMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Номер аудитории должен быть числом в диапазоне от 1 до 999</b>")
	return sb.String()
}

func subCreationErrorMessage(err error) string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Произошла ошибка при создании подписки:</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>%s</b>", err.Error()))
	return sb.String()
}

func subCreationSuccessMessage(labNumber, labAuditorium int) string {
	var sb strings.Builder
	sb.WriteString("<b>✅ Подписка создана!</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>📚 Лаба №%d</b>", labNumber))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>🚪 Аудитория №%d</b>", labAuditorium))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🔔 Вы получите уведомление, когда появится нужная запись</b>")
	return sb.String()
}

func unsubInvalidArgumentsMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Некорректные аргументы.</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>Использование: /unsub &lt;номер подписки в списке&gt;</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>Чтобы просмотреть список используйте команду /list</b>")
	return sb.String()
}

func unsubInvalidSubNumberMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Номер подписки должен быть числом</b>")
	return sb.String()
}

func unsubSubNumberValidationErrorMessage(subsLen int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>❌ Номер подписки должен быть числом в диапазоне от 1 до числа ваших подписок - %d</b>", subsLen))
	return sb.String()
}

func subsFetchingErrorMessage(err error) string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Произошла ошибка при получении списка подписок:</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>%s</b>", err.Error()))
	return sb.String()
}

func unsubErrorMessage(err error) string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Произошла ошибка при отписке:</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>%s</b>", err.Error()))
	return sb.String()
}

func unsubSuccessMessage(labNumber, labAuditorium int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("✅ Вы больше не подписаны на лабу №%d в ауд. №%d", labNumber, labAuditorium))
	return sb.String()
}

func listEmptySubsMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>🔍 У вас нет подписок на лабы.</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("Используйте команду /sub &lt;номер лабы&gt; &lt;номер аудитории&gt;")
	return sb.String()
}

func listSubsSuccessMessage(subs []model.Subscription) string {
	var sb strings.Builder
	sb.WriteString("<b>📋 Ваши подписки:</b>")
	sb.WriteString(repeatLineBreaks(2))
	for idx, sub := range subs {
		sb.WriteString(fmt.Sprintf("<b>%d. Лаба №%d, ауд. №%d</b>", idx+1, sub.LabNumber, sub.LabAuditorium))
		if idx == len(subs)-1 {
			break
		}
		sb.WriteString(repeatLineBreaks(2))
	}
	return sb.String()
}

func permissionDeniedErrorMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>❌ Доступ запрещён. Команда доступна только для разработчика</b>")
	return sb.String()
}

func statsSuccessMessage(snapshot *metrics.Metrics) string {
	uptime := time.Since(snapshot.StartTime)
	var sb strings.Builder
	sb.WriteString("<b>📊 Статистика сервиса</b>")
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🕐 Общее время работы:</b> ")
	sb.WriteString(formatDuration(uptime))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🔍 Опросы:</b>")
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Всего опросов: <b>%d</b>", snapshot.PollingMetrics.TotalPolls))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Режим: <b>%s</b>", formatPollingMode(snapshot.PollingMetrics.Mode)))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Ошибки парсинга: <b>%d</b>", snapshot.PollingMetrics.ParsingErrors))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Ошибки получения: <b>%d</b>", snapshot.PollingMetrics.FetchErrors))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Среднее время опроса: <b>%s</b>", snapshot.PollingMetrics.AveragePollingTime.Round(time.Millisecond)))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Среднее количество слотов: <b>%d</b>", snapshot.PollingMetrics.AverageSlotNumber))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🔔 Уведомления:</b>")
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Всего уведомлений: <b>%d</b>", snapshot.NotificationMetrics.TotalNotifications))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Размер кеша: <b>%d</b>", snapshot.NotificationMetrics.CacheLength))
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Среднее количество уведомлений: <b>%d</b>", snapshot.NotificationMetrics.AverageNotifications))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>📝 Подписки:</b>")
	sb.WriteString(repeatLineBreaks(1))
	sb.WriteString(fmt.Sprintf("  Активных подписок: <b>%d</b>", snapshot.SubscriptionMetrics.TotalSubscriptions))
	return sb.String()
}

func notifySuccessMessage(slot *model.Slot) string {
	var sb strings.Builder
	sb.WriteString("<b>🔥 Появилась запись!</b>")
	sb.WriteString(repeatLineBreaks(3))
	sb.WriteString(fmt.Sprintf("<b>📚 Лаба №%d. %s</b>", slot.LabNumber, slot.LabName))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString(fmt.Sprintf("<b>🚪 Аудитория №%d</b>", slot.LabAuditorium))
	sb.WriteString(repeatLineBreaks(2))
	sb.WriteString("<b>🗓️ Когда:</b>")
	sb.WriteString(repeatLineBreaks(2))
	for idx, dateTime := range slot.Available {
		sb.WriteString(fmt.Sprintf("<b>%d. %s</b>", idx+1, formatDateTime(dateTime)))
		sb.WriteString(repeatLineBreaks(2))
	}
	sb.WriteString(fmt.Sprintf("<b>🔗 <a href='%s'>Ссылка на запись</a></b>", slot.URL))
	return sb.String()
}

func repeatLineBreaks(breaks int) string {
	var sb strings.Builder
	for range breaks {
		sb.WriteString("\n")
	}
	return sb.String()
}
