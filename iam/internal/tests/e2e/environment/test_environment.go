//go:build integration

package environment

import (
	"context"
	"fmt"
)

// ClearState — очищает состояние Postgres и Redis между тестами
func (env *TestEnvironment) ClearState(ctx context.Context) error {
	_, err := env.Postgres.Pool().Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", usersTableName))
	if err != nil {
		return err
	}

	return env.Redis.Client().FlushDB(ctx).Err()
}

// DeleteUser — удаляет пользователя напрямую из Postgres, минуя API.
// Используется для воспроизведения ситуации "сессия существует, а пользователь — уже нет".
func (env *TestEnvironment) DeleteUser(ctx context.Context, userUUID string) error {
	_, err := env.Postgres.Pool().Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE uuid = $1", usersTableName), userUUID)
	return err
}
