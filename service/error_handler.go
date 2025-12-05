package service

import (
	"errors"

	"github.com/niiharamegumu/chronowork/internal/usecase"
	"github.com/rivo/tview"
)

// ErrorHandler handles error display in TUI.
type ErrorHandler struct {
	tui *TUI
}

// NewErrorHandler creates a new ErrorHandler.
func NewErrorHandler(tui *TUI) *ErrorHandler {
	return &ErrorHandler{tui: tui}
}

// ShowError displays an error message in a modal dialog.
// focusTarget is the widget name to focus after closing the modal.
func (h *ErrorHandler) ShowError(message string, focusTarget string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			h.tui.DeleteModal()
			h.tui.SetFocus(focusTarget)
		})
	h.tui.SetModal(modal)
	h.tui.SetFocus("modal")
}

// ShowErrorWithErr displays a user-friendly error message based on the error.
func (h *ErrorHandler) ShowErrorWithErr(err error, focusTarget string) {
	message := h.friendlyMessage(err)
	h.ShowError(message, focusTarget)
}

// friendlyMessage converts errors to user-friendly messages.
func (h *ErrorHandler) friendlyMessage(err error) string {
	if err == nil {
		return "不明なエラーが発生しました。"
	}

	// UseCaseErrorを型判定
	var ucErr *usecase.UseCaseError
	if errors.As(err, &ucErr) {
		return h.getMessageForCode(ucErr.Code)
	}

	// その他のエラー
	return "エラーが発生しました: " + err.Error()
}

// getMessageForCode returns a user-friendly message for the given error code.
func (h *ErrorHandler) getMessageForCode(code usecase.ErrorCode) string {
	messages := map[usecase.ErrorCode]string{
		usecase.ErrCodeNotFound:       "データが見つかりませんでした。",
		usecase.ErrCodeDuplicateToday: "この作業は今日既に作成されています。",
		usecase.ErrCodeValidation:     "入力内容に誤りがあります。",
		usecase.ErrCodePermission:     "権限がありません。",
	}

	if msg, ok := messages[code]; ok {
		return msg
	}
	return "エラーが発生しました。"
}
