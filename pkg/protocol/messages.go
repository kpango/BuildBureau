package protocol

// Task represents a unit of work assigned to an agent.
type Task struct {
	ID          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	AssignedBy  string `protobuf:"bytes,3,opt,name=assigned_by,json=assignedBy,proto3" json:"assigned_by,omitempty"`
}

type TaskResponse struct {
	TaskID   string `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Status   string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

type StatusUpdate struct {
	TaskID  string `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Status  string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	Result  string `protobuf:"bytes,4,opt,name=result,proto3" json:"result,omitempty"`
}

type StatusResponse struct {
	Received bool `protobuf:"varint,1,opt,name=received,proto3" json:"received,omitempty"`
}

type Message struct {
	From    string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	To      string `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	Content string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
}

type MessageResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}
