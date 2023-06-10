package job

import (
	"github.com/reugn/go-quartz/quartz"
)

type UpdateR2ImgJob struct {
	JobIns
}

func NewUpdateR2ImgJob() *UpdateR2ImgJob {
	userType := &UpdateR2ImgJob{}
	userType.IJob = interface{}(userType).(IJob)
	userType.Job = interface{}(userType).(quartz.Job)
	return userType
}

func (j *UpdateR2ImgJob) Description() string {
	return "Clear local images."
}

// not implemented
// func (j *UpdateR2ImgJob) ExecuteCore(logger *zap.SugaredLogger, cfg *helper.Config) (string, error) {
// folderPath := cfg.Prefix
// files, err := os.ReadDir(folderPath)
// if err != nil {

// 	return "", errors.Wrap(err, fmt.Sprintf("Error reading directory `%v`", folderPath))
// }

// for _, file := range files {
// 	if !file.IsDir() {
// 		p := folderPath + "/" + file.Name()
// 		err := os.Remove(p)
// 		if err != nil {
// 			return "", errors.Wrap(err, fmt.Sprintf("Error deleting file `%v`", p))
// 		}
// 	}
// }

// return fmt.Sprintf("removed %v image files", len(files)), nil
// }
