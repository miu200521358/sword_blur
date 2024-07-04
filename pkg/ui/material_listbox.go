package ui

import (
	"github.com/miu200521358/mlib_go/pkg/pmx"
	"github.com/miu200521358/walk/pkg/walk"
)

type MaterialListBox struct {
	*walk.ListBox
	MaterialListModel *MaterialListModel
}

func NewMaterialListBox(
	parent walk.Container,
) (*MaterialListBox, error) {
	lb, err := walk.NewListBox(parent)
	if err != nil {
		return nil, err
	}
	lb.SetMinMaxSize(walk.Size{Width: -1, Height: 100}, walk.Size{Width: -1, Height: 200})

	m := &MaterialListModel{
		items:                  make([]string, 0),
		itemsResetPublisher:    new(walk.EventPublisher),
		itemChangedPublisher:   new(walk.IntEventPublisher),
		itemsInsertedPublisher: new(walk.IntRangeEventPublisher),
		itemsRemovedPublisher:  new(walk.IntRangeEventPublisher),
	}
	lb.SetModel(m)

	return &MaterialListBox{ListBox: lb, MaterialListModel: m}, nil
}

func (lb *MaterialListBox) SetMaterials(
	materials *pmx.Materials, funcChange func(indexes []int),
) {
	for i := range materials.Len() {
		material := materials.Get(i)
		lb.MaterialListModel.items = append(lb.MaterialListModel.items, material.Name)
	}
	lb.MaterialListModel.PublishItemsReset()

	lb.SetSelectedIndexes(materials.GetIndexes())
	lb.SelectedIndexesChanged().Attach(func() {
		funcChange(lb.SelectedIndexes())
	})
}

type MaterialListModel struct {
	*walk.ReflectListModelBase
	itemsResetPublisher    *walk.EventPublisher
	itemChangedPublisher   *walk.IntEventPublisher
	itemsInsertedPublisher *walk.IntRangeEventPublisher
	itemsRemovedPublisher  *walk.IntRangeEventPublisher
	items                  []string
}

func (m *MaterialListModel) ItemCount() int {
	return len(m.items)
}

func (m *MaterialListModel) Value(index int) interface{} {
	return m.items[index]
}

func (m *MaterialListModel) Items() interface{} {
	return m.items
}

func (m *MaterialListModel) ItemsReset() *walk.Event {
	return m.itemsResetPublisher.Event()
}

func (m *MaterialListModel) ItemChanged() *walk.IntEvent {
	return m.itemChangedPublisher.Event()
}

func (m *MaterialListModel) ItemsInserted() *walk.IntRangeEvent {
	return m.itemsInsertedPublisher.Event()
}

func (m *MaterialListModel) ItemsRemoved() *walk.IntRangeEvent {
	return m.itemsRemovedPublisher.Event()
}

func (m *MaterialListModel) PublishItemsReset() {
	m.itemsResetPublisher.Publish()
}

func (m *MaterialListModel) PublishItemChanged(index int) {
	m.itemChangedPublisher.Publish(index)
}

func (m *MaterialListModel) PublishItemsInserted(from, to int) {
	m.itemsInsertedPublisher.Publish(from, to)
}

func (m *MaterialListModel) PublishItemsRemoved(from, to int) {
	m.itemsRemovedPublisher.Publish(from, to)
}
