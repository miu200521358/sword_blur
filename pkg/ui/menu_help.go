package ui

import (
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/walk/pkg/declarative"
)

func GetMenuItems() []declarative.MenuItem {
	return []declarative.MenuItem{
		declarative.Action{
			Text:        mi18n.T("概要"),
			OnTriggered: func() { mlog.ILT(mi18n.T("概要"), mi18n.T("概要メッセージ")) },
		},
		declarative.Separator{},
		declarative.Action{
			Text:        mi18n.T("Step1概要"),
			OnTriggered: func() { mlog.ILT(mi18n.T("Step1概要"), mi18n.T("Step1メッセージ")) },
		},
		declarative.Action{
			Text:        mi18n.T("Step2概要"),
			OnTriggered: func() { mlog.ILT(mi18n.T("Step2概要"), mi18n.T("Step2メッセージ")) },
		},
		declarative.Action{
			Text:        mi18n.T("Step3概要"),
			OnTriggered: func() { mlog.ILT(mi18n.T("Step3概要"), mi18n.T("Step3メッセージ")) },
		},
		declarative.Action{
			Text:        mi18n.T("Step4概要"),
			OnTriggered: func() { mlog.ILT(mi18n.T("Step4概要"), mi18n.T("Step4メッセージ")) },
		},
		declarative.Action{
			Text:        mi18n.T("Step5概要"),
			OnTriggered: func() { mlog.ILT(mi18n.T("Step5概要"), mi18n.T("Step5メッセージ")) },
		},
		declarative.Action{
			Text: mi18n.T("Step5プレビュー概要"),
			OnTriggered: func() {
				mlog.ILT(mi18n.T("Step5プレビュー概要"), mi18n.T("Step5プレビューメッセージ"))
			},
		},
	}
}
