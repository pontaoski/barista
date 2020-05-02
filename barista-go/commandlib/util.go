package commandlib

func ErrorEmbed(content string) Embed {
	return Embed{
		Title: EmbedHeader{
			Text: content,
		},
		Colour: 0xff0000,
	}
}
