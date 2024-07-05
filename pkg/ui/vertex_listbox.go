package ui

import (
	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"
)

type VertexListBox struct {
	*walk.ListBox
	VertexListModel *VertexListModel
}

func NewVertexListBox(
	parent walk.Container,
) (*VertexListBox, error) {
	// 複数選択リストボックス
	lb, err := walk.NewListBoxWithStyle(parent, win.LBS_MULTIPLESEL)
	if err != nil {
		return nil, err
	}
	lb.SetMinMaxSize(walk.Size{Width: -1, Height: 100}, walk.Size{Width: -1, Height: 200})

	m := &VertexListModel{
		items:                  make([]string, 0),
		itemsResetPublisher:    new(walk.EventPublisher),
		itemChangedPublisher:   new(walk.IntEventPublisher),
		itemsInsertedPublisher: new(walk.IntRangeEventPublisher),
		itemsRemovedPublisher:  new(walk.IntRangeEventPublisher),
	}
	lb.SetModel(m)

	return &VertexListBox{ListBox: lb, VertexListModel: m}, nil
}

type VertexListModel struct {
	*walk.ReflectListModelBase
	itemsResetPublisher    *walk.EventPublisher
	itemChangedPublisher   *walk.IntEventPublisher
	itemsInsertedPublisher *walk.IntRangeEventPublisher
	itemsRemovedPublisher  *walk.IntRangeEventPublisher
	items                  []string
}

func (m *VertexListModel) ItemCount() int {
	return len(m.items)
}

func (m *VertexListModel) Value(index int) interface{} {
	return m.items[index]
}

func (m *VertexListModel) Items() interface{} {
	return m.items
}

func (m *VertexListModel) ItemsReset() *walk.Event {
	return m.itemsResetPublisher.Event()
}

func (m *VertexListModel) ItemChanged() *walk.IntEvent {
	return m.itemChangedPublisher.Event()
}

func (m *VertexListModel) ItemsInserted() *walk.IntRangeEvent {
	return m.itemsInsertedPublisher.Event()
}

func (m *VertexListModel) ItemsRemoved() *walk.IntRangeEvent {
	return m.itemsRemovedPublisher.Event()
}

func (m *VertexListModel) PublishItemsReset() {
	m.itemsResetPublisher.Publish()
}

func (m *VertexListModel) PublishItemChanged(index int) {
	m.itemChangedPublisher.Publish(index)
}

func (m *VertexListModel) PublishItemsInserted(from, to int) {
	m.itemsInsertedPublisher.Publish(from, to)
}

func (m *VertexListModel) PublishItemsRemoved(from, to int) {
	m.itemsRemovedPublisher.Publish(from, to)
}
