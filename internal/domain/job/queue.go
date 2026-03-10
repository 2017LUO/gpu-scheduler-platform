package job

type QueueName string

const (
	QueueDefault   QueueName = "default"
	QueueInference QueueName = "inference"
	QueueTraining  QueueName = "training"
	QueueBatch     QueueName = "batch"
)
