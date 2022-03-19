package main

import (
	"crypto/rand"
	"flag"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"math/big"
	rand2 "math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func SendMessage(i *discordgo.InteractionCreate, Title, Description, Thumbnail string, Private ...bool) {
	if containsbool(Private, true) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       Title,
						Description: Description,
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: Thumbnail,
						},

						Color: 12189845,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Created by Leki#6796 | github.com/zLeki",
							IconURL: "http://www.leki.sbs/portfolio/img/image.png",
						},
					},
				},
			},
		})

	}else{
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{

				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       Title,
						Description: Description,
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: Thumbnail,
						},

						Color: 12189845,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Created by Leki#6796 | github.com/zLeki",
							IconURL: "http://www.leki.sbs/portfolio/img/image.png",
						},
					},
				},
			},
		})

	}
}
var s *discordgo.Session
func logmsg(i *discordgo.InteractionCreate, command string, args ...string) {
	s.ChannelMessageSendEmbed(logChannelID, &discordgo.MessageEmbed{
		Title:       "New Interaction",
		Description: i.Interaction.Member.User.Username+" ran a command.\n"+command+" "+strings.Join(args, " "),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: i.Interaction.Member.User.AvatarURL("80"),
		},

		Color: 12189845,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Created by Leki#6796 | github.com/zLeki",
			IconURL: "http://www.leki.sbs/portfolio/img/image.png",
		},
	})

}
func init() { flag.Parse() }
func init() {

	var err error
	s, err = discordgo.New("Bot ")
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}
func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := handler[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}
func main() {

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		err := s.UpdateStreamingStatus(0, s.State.User.Username+" | /help", "https://www.twitch.tv/amouranth")
		if err != nil {
			return
		}
	})
	s.Identify.Intents = discordgo.IntentsAll
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, s.State.Guilds[0].ID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, s.State.Guilds[0].ID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}


	log.Println("Gracefully shutdowning")
}
var (
	keys = make(map[string]struct{
		time string
		user *discordgo.User
		redeemed bool
	})
	admins = []string{
		"943281903765696612",
	}
	logChannelID = "935264163343781932"
	letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	roleID = "954426284455391282"
	commands = []*discordgo.ApplicationCommand{
	
		{
			Name: "whitelist",
			Description: "Whitelists a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name: "user",
					Description: "The user to whitelist",
					Type: discordgo.ApplicationCommandOptionUser,
					Required: true,
				},
			},

		},
		{
			Name: "unwhitelist",
			Description: "Un-whitelists a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name: "user",
					Description: "The user to blacklist",
					Type: discordgo.ApplicationCommandOptionUser,
					Required: true,
				},
			},
		},
		{
			Name: "createkey",
			Description: "Creates a key for a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name: "user",
					Description: "The user who can use the key",
					Type: discordgo.ApplicationCommandOptionUser,
					Required: true,
				},
			},
		},
		{
			Name: "redeem",
			Description: "Redeem a key",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name: "key",
					Description: "Enter key here.",
					Type: discordgo.ApplicationCommandOptionString,
					Required: true,
				},
			},
		},
		{
			Name: "check-key",
			Description: "Check a key for status",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name: "key",
					Description: "Enter key here.",
					Type: discordgo.ApplicationCommandOptionString,
					Required: true,
				},
			},
		},
		{
			Name: "blacklist",
			Description: "Blacklist a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name: "user",
					Description: "The user to blacklist.",
					Type: discordgo.ApplicationCommandOptionUser,
					Required: true,
				},
			},
		},
		{
			Name: "unblacklist",
			Description: "Pardon a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name: "user",
					Description: "The user to unblacklist.",
					Type: discordgo.ApplicationCommandOptionUser,
					Required: true,
				},
			},
		},
	}
	handler = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		
		"unblacklist": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			userID := i.ApplicationCommandData().Options[0].UserValue(s)
			f, err := os.OpenFile("./data/pardoned.log",
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			defer f.Close()
			if _, err := f.WriteString(userID.ID+"\n"); err != nil {
				log.Println(err)
				return
			}
			SendMessage(i, "Unblacklisted", "Unblacklisted the user, "+i.Member.User.Username, "https://i.imgur.com/I5ttBFi.png")

		},
		"check-key": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			key := i.ApplicationCommandData().Options[0].StringValue()
			if keys[key].user == nil {
				SendMessage(i, "Warning", "Key does not exist", "https://i.imgur.com/NgxYShD.png")
				logmsg(i, "/check-key", key);return
			}
			SendMessage(i, "Key data.", "`Redeemed:` **"+strconv.FormatBool(keys[key].redeemed)+"**\n`Time redeemed/created:` **"+keys[key].time+"**\n`Redeemed/Created by:` **"+keys[key].user.Username+"**", "https://i.imgur.com/NldSwaZ.png", true)
		},
		"redeem": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			data, _ := ioutil.ReadFile("./data/text.log")
			users := strings.Split(string(data), "\n")
			if contains(users, i.Member.User.ID) {
				data2, _ := ioutil.ReadFile("./data/text.log")
				users2 := strings.Split(string(data2), "\n")
				if !contains(users2, i.Member.User.ID) {
					SendMessage(i, "You are blacklisted sorry.", "403 Denied", "https://i.imgur.com/qs4QOjF.png")
				}
			}
			key := i.ApplicationCommandData().Options[0].StringValue()
			if keys[key].redeemed == true {
				SendMessage(i, "Key used.", "Key was already used by "+keys[key].user.Username+" at "+keys[key].time+". If this is a mistake contact support.", "https://i.imgur.com/qs4QOjF.png",true)
				logmsg(i, "/redeem", key);return
			}
			log.Println(keys[key].user.ID, i.Interaction.Member.User.ID)
			if keys[key].user.ID == i.Interaction.Member.User.ID {
				keys[key] = struct {
					time     string
					user     *discordgo.User
					redeemed bool
				}{time: time.Now().String(), user: i.Member.User, redeemed: true}
				SendMessage(i, "Key accepted.", "You are now whitelisted.", "https://i.imgur.com/I5ttBFi.png", true)
				s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, roleID)
			}else{
				SendMessage(i, "Access denied.", "This key is reserved for, **"+keys[key].user.Username+"**", "https://i.imgur.com/qs4QOjF.png")
			}

		},
		"whitelist": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if CheckPermissions(i, i.Member.User.ID) {
					userID := i.ApplicationCommandData().Options[0].UserValue(s)
					member, _ := s.GuildMember(i.GuildID, userID.ID)
					privateChannel, _ := s.UserChannelCreate(userID.ID)
					for _,v := range member.Roles {
						if v == roleID {
							SendMessage(i, "Warning", "User is already whitelisted", "https://i.imgur.com/NgxYShD.png")
							logmsg(i, "/whitelist", userID.Username);return
						}
					}
					s.GuildMemberRoleAdd(i.GuildID, userID.ID, roleID)
					SendMessage(i, "Success", "User whitelisted", "https://i.imgur.com/I5ttBFi.png")
					s.ChannelMessageSendEmbed(privateChannel.ID, &discordgo.MessageEmbed{
						Title: "Whitelisted",
						Description: "You have been whitelisted.",
						Color: 12189845,
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: "https://i.imgur.com/I5ttBFi.png",
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text: "Created by Leki#6796 | github.com/zLeki",
							IconURL: "http://www.leki.sbs/portfolio/img/image.png",
						},
					})
					logmsg(i, "/whitelist", userID.Username);return
				}

			SendMessage(i, "Access Denied", "You do not have access to this command. Sorry bout that", "https://i.imgur.com/qs4QOjF.png", true)

		},
		"blacklist": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			userID := i.ApplicationCommandData().Options[0].UserValue(s)
			if CheckPermissions(i, i.Member.User.ID) {
				f, err := os.OpenFile("./data/text.log",
					os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Println(err)
				}
				defer f.Close()
				if _, err := f.WriteString(userID.ID+"\n"); err != nil {
					log.Println(err)
					return
				}
				SendMessage(i, "Success", "Successfully blacklisted, "+userID.Username, "https://i.imgur.com/I5ttBFi.png")
			}
		},
		"unwhitelist": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			userID := i.ApplicationCommandData().Options[0].UserValue(s)
			member, _ := s.GuildMember(i.GuildID, userID.ID)
			privateChannel, _ := s.UserChannelCreate(userID.ID)
			if CheckPermissions(i, i.Member.User.ID) {
				for _, v := range member.Roles {
					if v == roleID {
						s.GuildMemberRoleRemove(i.GuildID, userID.ID, roleID)
						SendMessage(i, "Success", "User unwhitelisted", "https://i.imgur.com/I5ttBFi.png")
						s.ChannelMessageSendEmbed(privateChannel.ID, &discordgo.MessageEmbed{
							Title:       "Whitelisted",
							Description: "You have been un-whitelisted.",
							Color:       12189845,
							Thumbnail: &discordgo.MessageEmbedThumbnail{
								URL: "https://i.imgur.com/I5ttBFi.png",
							},
							Footer: &discordgo.MessageEmbedFooter{
								Text:    "Created by Leki#6796 | github.com/zLeki",
								IconURL: "http://www.leki.sbs/portfolio/img/image.png",
							},
						})
						logmsg(i, "/unwhitelist", userID.Username);return
					}
				}
				SendMessage(i, "Warning", "User is not whitelisted.", "https://i.imgur.com/NgxYShD.png")
				logmsg(i, "/unwhitelist", userID.Username);return
			}
			SendMessage(i, "Access Denied", "You do not have access to this command. Sorry bout that", "https://i.imgur.com/qs4QOjF.png", true)

		},
		"createkey": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if CheckPermissions(i, i.Member.User.ID) {
				ret := make([]byte, 16)
				for i := 0; i < 16; i++ {
					num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
					if err != nil {
						panic(err)
					}
					ret[i] = letters[num.Int64()]
				}
				keys[string(ret)] = struct {
					time     string
					user     *discordgo.User
					redeemed bool
				}{time: time.Now().String(), user: i.ApplicationCommandData().Options[0].UserValue(s), redeemed: false}
				SendMessage(i, "Here is your private key.", "Only has one use, and it expires in a few hours.\n`"+string(ret)+"`", "https://i.imgur.com/NldSwaZ.png", true)
				logmsg(i, "/createkey", i.Member.User.Username);return
			}
			SendMessage(i, "Access Denied", "You do not have access to this command. Sorry bout that", "https://i.imgur.com/qs4QOjF.png", true)

		},

	}
)
func contains(a []string, b string) bool {
	for _,v := range a  {
		if v == b {
			return true
		}
	}
	return false
}
func containsbool(a []bool, b bool) bool {
	for _,v := range a  {
		if v == b {
			return true
		}
	}
	return false
}
func CheckPermissions(i *discordgo.InteractionCreate, userid string) bool {
	selfUser, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		log.Printf("Cannot get self user: %v", err)
		return false
	}
	for _,v := range selfUser.Roles {
		if contains(admins, v) {
			return true
		}
	}
	return false
}
//Images: https://i.imgur.com/v2n7qPs.png-Ping, https://i.imgur.com/NldSwaZ.png-Info, https://i.imgur.com/qs4QOjF.png-Error, https://i.imgur.com/I5ttBFi.png-Success, https://i.imgur.com/NgxYShD.png-Warning
