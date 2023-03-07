package domain

import (
	"errors"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

type User struct {
	Id      string
	Email   dp.Email
	Account dp.Account
}

// SoftwarePkgReviewComment
type SoftwarePkgReviewComment struct {
	Id        string
	CreatedAt int64
	Author    dp.Account
	Content   dp.ReviewComment
}

func NewSoftwarePkgReviewComment(
	author dp.Account, content dp.ReviewComment,
) SoftwarePkgReviewComment {
	return SoftwarePkgReviewComment{
		CreatedAt: utils.Now(),
		Author:    author,
		Content:   content,
	}
}

// SoftwarePkgApplication
type SoftwarePkgApplication struct {
	SourceCode        SoftwarePkgSourceCode
	PackageDesc       dp.PackageDesc
	PackagePlatform   dp.PackagePlatform
	ImportingPkgSig   dp.ImportingPkgSig
	ReasonToImportPkg dp.ReasonToImportPkg
}

type SoftwarePkgSourceCode struct {
	Address dp.URL
	License dp.License
}

// SoftwarePkgBasicInfo
type SoftwarePkgBasicInfo struct {
	Id           string
	PkgName      dp.PackageName
	Importer     dp.Account
	RepoLink     dp.URL
	Phase        dp.PackagePhase
	Frozen       bool
	ReviewResult dp.PackageReviewResult
	AppliedAt    int64
	Application  SoftwarePkgApplication
	ApprovedBy   []dp.Account
	RejectedBy   []dp.Account
	RelevantPR   dp.URL
}

func (entity *SoftwarePkgBasicInfo) CanAddReviewComment() bool {
	return entity.Phase.IsReviewing() || entity.Phase.IsCreatingRepo()
}

// change the status of "creating repo"
// send out the event
// notify the importer
func (entity *SoftwarePkgBasicInfo) ApproveBy(user dp.Account) (bool, error) {
	if !entity.Phase.IsReviewing() || entity.Frozen || entity.RelevantPR == nil {
		return false, errors.New("not ready")
	}

	entity.ApprovedBy = append(entity.ApprovedBy, user)

	approved := false
	// only set the result once to avoid duplicate case.
	if len(entity.ApprovedBy) == 2 {
		entity.ReviewResult = dp.PkgReviewResultApproved
		entity.Phase = dp.PackagePhaseCreatingRepo
		approved = true
	}

	return approved, nil
}

// notify the importer
func (entity *SoftwarePkgBasicInfo) RejectBy(user dp.Account) (bool, error) {
	if !entity.Phase.IsReviewing() {
		return false, errors.New("can't do this")
	}

	entity.RejectedBy = append(entity.RejectedBy, user)

	rejected := false
	// only set the result once to avoid duplicate case.
	if len(entity.RejectedBy) == 1 {
		entity.ReviewResult = dp.PkgReviewResultRejected
		entity.Phase = dp.PackagePhaseClosed
		rejected = true
	}

	return rejected, nil
}

func (entity *SoftwarePkgBasicInfo) Abandon(user dp.Account) error {
	if !entity.Phase.IsReviewing() {
		return errors.New("can't do this")
	}

	if !dp.IsSameAccount(user, entity.Importer) {
		return errors.New("not the importer")
	}

	entity.Phase = dp.PackagePhaseClosed

	return nil
}

// SoftwarePkg
type SoftwarePkg struct {
	SoftwarePkgBasicInfo

	Comments []SoftwarePkgReviewComment
}

func NewSoftwarePkg(user dp.Account, name dp.PackageName, app *SoftwarePkgApplication) SoftwarePkgBasicInfo {
	return SoftwarePkgBasicInfo{
		PkgName:     name,
		Importer:    user,
		Phase:       dp.PackagePhaseReviewing,
		Frozen:      true,
		Application: *app,
		AppliedAt:   utils.Now(),
	}
}
