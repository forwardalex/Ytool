package wrokpool

type Jobs interface {
	GetJob()
	Worker()
}

type innerJob struct {
}

func (j *innerJob) GetJob() {

}
func (j *innerJob) Worker() {

}
