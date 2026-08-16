package domain

import (
	"errors"
	"time"
)

type RequestType string

const (
	RequestTypeMaterial     RequestType = "material"
	RequestTypeWork         RequestType = "work"
	RequestTypeStatusUpdate RequestType = "status_update"
)

type RequestStatus string

const (
	StatusNew        RequestStatus = "new"
	StatusInProgress RequestStatus = "in_progress"
	StatusDone       RequestStatus = "done"
	StatusRejected   RequestStatus = "rejected"
	StatusCancelled  RequestStatus = "cancelled"
)

var ErrInvalidStatusTransition = errors.New("invalid status transition")

// Request — агрегат. Сам знает, какие переходы статуса для него допустимы:
// бизнес-правило живёт в домене, а не в сервисе и не в хендлере.
type Request struct {
	ID          int64         // Первичный ключ (PK) в БД, генерируется при INSERT
	Object      string        // Строковый идентификатор объекта (например, "Объект в кп Мечта-2")
	Type        RequestType   // Тип заявки (enum, защищает от невалидных значений)
	Title       string        // Короткий заголовок для отображения в списках
	Description string        // Полное текстовое описание заявки (опционально)
	Status      RequestStatus // Текущий статус.!!! Менять только через TransitionTo() !!!
	CreatedBy   int64         // ID пользователя-создателя (FK) для проверки прав
	CreatedAt   time.Time     // Время создания (устанавливается БД, не меняется)
	UpdatedAt   time.Time     // Время последнего обновления статуса
}

// CanTransitionTo - чистая функция без побочных эффектов: только читает
// текущий статус, ничего не меняет.
func (r *Request) CanTransitionTo(newStatus RequestStatus) bool {
	switch r.Status {
	case StatusNew:
		return newStatus == StatusInProgress || newStatus == StatusRejected || newStatus == StatusCancelled
	case StatusInProgress:
		return newStatus == StatusDone || newStatus == StatusRejected || newStatus == StatusCancelled
	case StatusDone, StatusRejected, StatusCancelled:
		//терминальные статусы — переходов нет, намеренно не в default, для читаемости
		return false
	default:
		return false
	}
}

func (r *Request) TransitionTo(newStatus RequestStatus) error {
	if !r.CanTransitionTo(newStatus) {
		return ErrInvalidStatusTransition
	}

	r.Status = newStatus
	r.UpdatedAt = time.Now()
	return nil
}

// StatusHistoryEntry - запись аудита. OldStatus/NewStatus не nullable:
// запись создаётся только на реальном переходе (через TransitionTo),
// а не в момент создания заявки - тот факт виден по Request.CreatedAt.
type StatusHistoryEntry struct {
	ID        int64         // Уникальный ID записи в таблице истории
	RequestID int64         // К какой заявке относится (Foreign Key)
	OldStatus RequestStatus // Какой был статус ДО изменения
	NewStatus RequestStatus // Какой стал статус ПОСЛЕ изменения
	Comment   string        // Комментарий: почему изменили
	ChangedBy int64         // Кто изменил (ID пользователя)
	ChangedAt time.Time     // Когда изменили (временная метка)
}
