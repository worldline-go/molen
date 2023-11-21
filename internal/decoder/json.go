package decoder

import "encoding/json"

type Messages []json.RawMessage

func (m *Messages) UnmarshalJSON(b []byte) error {
	type messages Messages
	msgs := messages{}

	err := json.Unmarshal(b, &msgs)
	if err != nil {
		message := json.RawMessage{}

		err = json.Unmarshal(b, &message)
		if err != nil {
			return err
		}

		*m = append(*m, message)

		return nil
	} else {
		*m = Messages(msgs)
	}

	return nil
}
