package ssoclient

import (
	"github.com/yourusername/sso-client-go/pkg/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	primaryDB   *gorm.DB
	secondaryDB *gorm.DB
}

func NewUserRepository(primaryDB *gorm.DB, secondaryDB *gorm.DB) *UserRepository {
	return &UserRepository{
		primaryDB:   primaryDB,
		secondaryDB: secondaryDB,
	}
}

// tryDBs attempts to execute the given function first on primary DB, then on secondary if primary fails
func (r *UserRepository) tryDBs(operation func(*gorm.DB) error) error {
	if err := operation(r.primaryDB); err != nil && r.secondaryDB != nil {
		return operation(r.secondaryDB)
	}
	return nil
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	var lastErr error

	// Try primary DB
	result := r.primaryDB.First(&user, id)
	if result.Error == nil {
		return &user, nil
	}
	lastErr = result.Error

	// Try secondary DB if available
	if r.secondaryDB != nil {
		result = r.secondaryDB.First(&user, id)
		if result.Error == nil {
			return &user, nil
		}
		lastErr = result.Error
	}

	return nil, lastErr
}

func (r *UserRepository) Create(user *models.User) error {
	return r.tryDBs(func(db *gorm.DB) error {
		return db.Create(user).Error
	})
}

func (r *UserRepository) Update(user *models.User) error {
	return r.tryDBs(func(db *gorm.DB) error {
		return db.Save(user).Error
	})
}

func (r *UserRepository) Delete(id uint) error {
	return r.tryDBs(func(db *gorm.DB) error {
		return db.Delete(&models.User{}, id).Error
	})
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	var lastErr error

	// Try primary DB
	result := r.primaryDB.Where("email = ?", email).First(&user)
	if result.Error == nil {
		return &user, nil
	}
	lastErr = result.Error

	// Try secondary DB if available
	if r.secondaryDB != nil {
		result = r.secondaryDB.Where("email = ?", email).First(&user)
		if result.Error == nil {
			return &user, nil
		}
		lastErr = result.Error
	}

	return nil, lastErr
}

func (r *UserRepository) FindByJTI(jti string) (uint, error) {
	var token models.UserAccessToken
	var lastErr error

	// Try primary DB
	result := r.primaryDB.Where("jti = ?", jti).First(&token)
	if result.Error == nil {
		return token.UserID, nil
	}
	lastErr = result.Error

	// Try secondary DB if available
	if r.secondaryDB != nil {
		result = r.secondaryDB.Where("jti = ?", jti).First(&token)
		if result.Error == nil {
			return token.UserID, nil
		}
		lastErr = result.Error
	}

	return 0, lastErr
}

func (r *UserRepository) GetLastSshKey() (*models.SshKey, error) {
	var sshKey models.SshKey
	var lastErr error

	// Try primary DB
	result := r.primaryDB.Order("id desc").First(&sshKey)
	if result.Error == nil {
		return &sshKey, nil
	}
	lastErr = result.Error

	// Try secondary DB if available
	if r.secondaryDB != nil {
		result = r.secondaryDB.Order("id desc").First(&sshKey)
		if result.Error == nil {
			return &sshKey, nil
		}
		lastErr = result.Error
	}

	return nil, lastErr
}

func (r *UserRepository) CreateAccessToken(userID uint, jti string) error {
	token := &models.UserAccessToken{
		UserID: userID,
		JTI:    jti,
	}
	return r.tryDBs(func(db *gorm.DB) error {
		return db.Create(token).Error
	})
}

func (r *UserRepository) DeleteAccessToken(jti string) error {
	return r.tryDBs(func(db *gorm.DB) error {
		return db.Where("jti = ?", jti).Delete(&models.UserAccessToken{}).Error
	})
}

// Helper methods for SSH keys

func (r *UserRepository) CreateSshKey(key string) error {
	sshKey := &models.SshKey{
		PrivateRsaKey: key,
	}
	return r.tryDBs(func(db *gorm.DB) error {
		return db.Create(sshKey).Error
	})
}

func (r *UserRepository) DeleteSshKey(id uint) error {
	return r.tryDBs(func(db *gorm.DB) error {
		return db.Delete(&models.SshKey{}, id).Error
	})
}
