package views

type GoalCancelForm struct {
	Form
	Destination string
}

func NewGoalCancelForm(balance int) *GoalCancelForm {
	f := GoalCancelForm{}
	f.Initial = true
	return &f
}

func (v *GoalCancelForm) ValidateDestination() (msgs []string) {
	if v.Initial {
		return
	}

	if v.Destination == "" {
		msgs = append(msgs, "Destination is required")
	}

	return msgs
}

func (v *GoalCancelForm) DestinationHasError() bool {
	return len(v.ValidateDestination()) > 0
}

func (v *GoalCancelForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateDestination()...)
	return msgs
}
