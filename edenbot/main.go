package edenbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yamamushi/EscapingEden/edenconfig"
	"github.com/yamamushi/EscapingEden/edendb"
	"github.com/yamamushi/EscapingEden/logging"
	"github.com/yamamushi/EscapingEden/messages"
)

type EdenBot struct {
	DB     edendb.DatabaseType
	Log    logging.LoggerType
	Config *edenconfig.Config
	dg     *discordgo.Session

	InputChannel  chan messages.EdenbotMessage
	OutputChannel chan messages.SystemManagerMessage
}

func NewEdenBot(input chan messages.EdenbotMessage, output chan messages.SystemManagerMessage, db edendb.DatabaseType,
	logger logging.LoggerType, conf *edenconfig.Config) *EdenBot {
	return &EdenBot{
		DB:            db,
		Log:           logger,
		Config:        conf,
		InputChannel:  input,
		OutputChannel: output,
	}
}

func (eb *EdenBot) Init() error {
	return nil //TODO
}

func (eb *EdenBot) Run(startNotify chan bool) error {
	// First we start the discord bot using the credentials from eb.config
	// Then we start the manager service
	eb.Log.Println(logging.LogInfo, "Edenbot connecting to discord...")
	dg, err := discordgo.New("Bot " + eb.Config.Discord.Token)
	if err != nil {
		return err
	}
	eb.dg = dg
	err = eb.dg.Open()
	if err != nil {
		return err
	}
	eb.dg.State.TrackMembers = true
	eb.dg.State.TrackChannels = true
	eb.dg.State.TrackEmojis = true
	eb.dg.State.TrackPresences = true
	eb.dg.State.TrackRoles = true
	eb.dg.State.TrackVoice = true
	if err != nil {
		eb.Log.Println(logging.LogError, "Error opening Discord session: ", err)
		return err
	}
	eb.Log.Println(logging.LogInfo, "Discord Connection Established")
	err = eb.dg.UpdateGameStatus(0, "Eden is online.")
	if err != nil {
		return err
	}

	eb.Log.Println(logging.LogInfo, "Edenbot starting manager service...")
	go eb.HandleInput(startNotify)
	eb.dg.AddHandler(eb.HandleVerify)
	//eb.dg.ChannelMessageSend(eb.Config.Discord.RegistrationChannelID, "Edenbot is online.")

	return nil
}

func (eb *EdenBot) Shutdown() error {
	err := eb.dg.Close()
	if err != nil {
		return err
	}
	return nil
}
