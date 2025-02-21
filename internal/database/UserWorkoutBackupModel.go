package database

import "github.com/jmoiron/sqlx"

type UserWorkoutBackup struct {
	UserID     int64  `json:"userId"`
	BackupPath string `json:"backupPath"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type UserWorkoutBackupModel struct {
	*sqlx.DB
}

func (model *UserWorkoutBackupModel) GetByUserId(userID int64) (UserWorkoutBackup, error) {
	u := new(UserWorkoutBackup)

	row := model.DB.QueryRow("SELECT user_id, backup_path, created_at, updated_at FROM user_workout_backups WHERE user_id = ?", userID)

	err := row.Scan(&u.UserID, &u.BackupPath, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return UserWorkoutBackup{}, err
	}

	return *u, nil
}

func (model *UserWorkoutBackupModel) CreateWorkoutBackup(userID int64, backupPath string) error {
	_, err := model.DB.Exec("INSERT INTO user_workout_backups (user_id, backup_path) VALUES (?, ?)", userID, backupPath)

	return err
}

func (model *UserWorkoutBackupModel) TouchWorkoutBackup(uwb *UserWorkoutBackup) error {
	// Just update the updated_at

	_, err := model.DB.Exec("UPDATE user_workout_backups SET updated_at = NOW() WHERE user_id = ? AND backup_path = ?", uwb.UserID, uwb.BackupPath)

	return err
}
