package bullet_engine

import "vixac.com/got/engine"

type MockSummaryStore struct {
	errorToThrow error

	upsertId  *engine.SummaryId
	upsertAgg *engine.Summary

	aggs     map[engine.SummaryId]engine.Summary
	fetchIds []engine.SummaryId
}

func MakeMockSummaryStore() MockSummaryStore {
	aggs := make(map[engine.SummaryId]engine.Summary)
	return MockSummaryStore{
		aggs: aggs,
	}
}
func (m *MockSummaryStore) UpsertAggregate(id engine.SummaryId, agg engine.Summary) error {
	m.upsertAgg = &agg
	m.upsertId = &id
	return m.errorToThrow
}

func (m *MockSummaryStore) UpsertManyAggregates(aggs map[engine.SummaryId]engine.Summary) error {

	for k, v := range aggs {
		m.aggs[k] = v
	}
	return m.errorToThrow
}
func (m *MockSummaryStore) Fetch(ids []engine.SummaryId) (map[engine.SummaryId]engine.Summary, error) {
	found := make(map[engine.SummaryId]engine.Summary)
	for _, id := range ids {
		existing, ok := m.aggs[id]
		if ok {
			found[id] = existing
		}
	}
	m.fetchIds = ids
	return found, m.errorToThrow
}
func (m *MockSummaryStore) Delete(ids []engine.SummaryId) error {
	return m.errorToThrow
}
