package components

import (
	"restclient/internal/tui/styles"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lv2 "charm.land/lipgloss/v2"
)

type AuthType int

const (
	AuthNone AuthType = iota
	AuthBearer
	AuthBasic
)

var authTypeLabels = []string{"None", "Bearer Token", "Basic Auth"}

// AuthEditor lets the user pick an authorization type and fill in the
// corresponding credentials.
type AuthEditor struct {
	Type AuthType

	Token    textinput.Model
	Username textinput.Model
	Password textinput.Model

	width int
}

func NewAuthEditor() AuthEditor {
	token := textinput.New()
	token.Placeholder = "token"
	token.SetVirtualCursor(true)

	username := textinput.New()
	username.Placeholder = "username"
	username.SetVirtualCursor(true)

	password := textinput.New()
	password.Placeholder = "password"
	password.EchoMode = textinput.EchoPassword
	password.SetVirtualCursor(true)

	return AuthEditor{Token: token, Username: username, Password: password}
}

func (m AuthEditor) SetWidth(width int) AuthEditor {
	m.width = width
	fieldWidth := width - 12
	if fieldWidth < 10 {
		fieldWidth = 10
	}
	m.Token.SetWidth(fieldWidth)
	m.Username.SetWidth(fieldWidth)
	m.Password.SetWidth(fieldWidth)
	return m
}

func (m AuthEditor) blurAll() AuthEditor {
	m.Token.Blur()
	m.Username.Blur()
	m.Password.Blur()
	return m
}

func (m AuthEditor) typeLabelView(t AuthType) string {
	style := lv2.NewStyle().Foreground(styles.TextSubtle).Padding(0, 1)
	if m.Type == t {
		style = style.Foreground(styles.TextNormal).Bold(true).Underline(true)
	}
	return style.Render(authTypeLabels[t])
}

func (m AuthEditor) typeRowView() string {
	parts := make([]string, 0, len(authTypeLabels)*2-1)
	for i := range authTypeLabels {
		if i > 0 {
			parts = append(parts, " ")
		}
		parts = append(parts, m.typeLabelView(AuthType(i)))
	}
	return lv2.JoinHorizontal(lv2.Top, parts...)
}

// typeBounds returns the start/end column of each selectable auth type
// label, reusing the tabBounds type shared with TabsModel.
func (m AuthEditor) typeBounds() []tabBounds {
	bounds := make([]tabBounds, len(authTypeLabels))
	cursor := 0
	for i := range authTypeLabels {
		w := lv2.Width(m.typeLabelView(AuthType(i)))
		bounds[i] = tabBounds{start: cursor, end: cursor + w}
		cursor += w + 1
	}
	return bounds
}

// Height reports how many terminal rows this editor currently occupies.
func (m AuthEditor) Height() int {
	switch m.Type {
	case AuthBearer:
		return 3 // type row, gap, token row
	case AuthBasic:
		return 4 // type row, gap, username row, password row
	default:
		return 3
	}
}

// HandleClick processes a click at coordinates relative to the top-left
// of this editor.
func (m AuthEditor) HandleClick(relX, relY int) AuthEditor {
	if relY == 0 {
		for i, b := range m.typeBounds() {
			if relX >= b.start && relX < b.end {
				m.Type = AuthType(i)
				return m.blurAll()
			}
		}
		return m.blurAll()
	}

	m = m.blurAll()

	switch m.Type {
	case AuthBearer:
		if relY == 2 {
			m.Token.Focus()
		}
	case AuthBasic:
		if relY == 2 {
			m.Username.Focus()
		} else if relY == 3 {
			m.Password.Focus()
		}
	}

	return m
}

func (m AuthEditor) Update(msg tea.Msg) (AuthEditor, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.Token, cmd = m.Token.Update(msg)
	cmds = append(cmds, cmd)
	m.Username, cmd = m.Username.Update(msg)
	cmds = append(cmds, cmd)
	m.Password, cmd = m.Password.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m AuthEditor) View() string {
	typeRow := m.typeRowView()
	gap := lv2.NewStyle().Height(1).Render("")

	var fields string
	switch m.Type {
	case AuthBearer:
		fields = lv2.JoinHorizontal(lv2.Top,
			lv2.NewStyle().Foreground(styles.TextSubtle).Width(10).Render("Token"),
			m.Token.View())
	case AuthBasic:
		userRow := lv2.JoinHorizontal(lv2.Top,
			lv2.NewStyle().Foreground(styles.TextSubtle).Width(10).Render("Username"),
			m.Username.View())
		passRow := lv2.JoinHorizontal(lv2.Top,
			lv2.NewStyle().Foreground(styles.TextSubtle).Width(10).Render("Password"),
			m.Password.View())
		fields = lv2.JoinVertical(lv2.Left, userRow, passRow)
	default:
		fields = lv2.NewStyle().
			Foreground(styles.TextMuted).
			Render("This request does not use any authorization")
	}

	return lv2.JoinVertical(lv2.Left, typeRow, gap, fields)
}
