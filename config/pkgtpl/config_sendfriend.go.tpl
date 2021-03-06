package sendfriend

var PackageConfiguration = config.NewConfiguration(
	&config.Section{
		ID:        "sendfriend",
		Label:     "Email to a Friend",
		SortOrder: 120,
		Scope:     config.ScopePermAll,
		Groups: config.GroupSlice{
			&config.Group{
				ID:        "email",
				Label:     `Email Templates`,
				Comment:   ``,
				SortOrder: 1,
				Scope:     config.ScopePermAll,
				Fields: config.FieldSlice{
					&config.Field{
						// Path: `sendfriend/email/enabled`,
						ID:           "enabled",
						Label:        `Enabled`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    1,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      true,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sendfriend/email/template`,
						ID:           "template",
						Label:        `Select Email Template`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    2,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      `sendfriend_email_template`,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Email\Template
					},

					&config.Field{
						// Path: `sendfriend/email/allow_guest`,
						ID:           "allow_guest",
						Label:        `Allow for Guests`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    3,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Config\Model\Config\Source\Yesno
					},

					&config.Field{
						// Path: `sendfriend/email/max_recipients`,
						ID:           "max_recipients",
						Label:        `Max Recipients`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    4,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      5,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sendfriend/email/max_per_hour`,
						ID:           "max_per_hour",
						Label:        `Max Products Sent in 1 Hour`,
						Comment:      ``,
						Type:         config.TypeText,
						SortOrder:    5,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      5,
						BackendModel: nil,
						SourceModel:  nil,
					},

					&config.Field{
						// Path: `sendfriend/email/check_by`,
						ID:           "check_by",
						Label:        `Limit Sending By`,
						Comment:      ``,
						Type:         config.TypeSelect,
						SortOrder:    6,
						Visible:      config.VisibleYes,
						Scope:        config.ScopePermAll,
						Default:      false,
						BackendModel: nil,
						SourceModel:  nil, // Magento\Sendfriend\Model\Source\Checktype
					},
				},
			},
		},
	},
)
