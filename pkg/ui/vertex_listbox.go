package ui

import (
	"slices"

	"github.com/miu200521358/mlib_go/pkg/mutils"
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

func (lb *VertexListBox) Clear() {
	lb.VertexListModel.items = make([]string, 0)
	lb.VertexListModel.PublishItemsReset()
}

func (lb *VertexListBox) GetItemValues() []int {
	items := make([]int, 0)
	for _, item := range lb.VertexListModel.items {
		// カンマ繋ぎの文字列を数値リストに変換
		itemValues, err := mutils.SplitCommaSeparatedInts(item)
		if err == nil {
			items = append(items, itemValues...)
		}
	}
	return items
}

func (lb *VertexListBox) SetItem(items []int) {
	// 順不同なのでとりあえずソート
	slices.Sort(items)
	// 頂点番号リストをカンマで繋ぐ
	itemStr := mutils.JoinIntsWithComma(items)
	isRemoved := false
	for i, it := range lb.VertexListModel.items {
		if it == itemStr {
			lb.RemoveItem(i)
			isRemoved = true
			break
		}
	}
	if !isRemoved {
		lb.AppendItem(itemStr)
	}
}

func (lb *VertexListBox) AppendItem(item string) {
	lb.VertexListModel.items = append(lb.VertexListModel.items, item)
	lb.VertexListModel.PublishItemsInserted(len(lb.VertexListModel.items)-1, len(lb.VertexListModel.items)-1)
}

func (lb *VertexListBox) RemoveItem(index int) {
	lb.VertexListModel.items = append(lb.VertexListModel.items[:index], lb.VertexListModel.items[index+1:]...)
	lb.VertexListModel.PublishItemsRemoved(index, index)
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
