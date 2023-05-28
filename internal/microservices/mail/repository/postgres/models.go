package postgres

type User2Fake struct {
	UserID uint64
	FakeID uint64
}

func (User2Fake) TableName(schemaName string) string {
	return schemaName + ".user2fake"
}
