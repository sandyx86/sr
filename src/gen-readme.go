package main

import (
	"encoding/xml"
	"log"
	"os"
)

type XMLRun struct {
	XMLName xml.Name `xml:"Run"`
	//Text         string   `xml:",chardata"`
	Version      string `xml:"version,attr"`
	GameIcon     string `xml:"GameIcon"`
	GameName     string `xml:"GameName"`
	CategoryName string `xml:"CategoryName"`
	Metadata     struct {
		//Text string `xml:",chardata"`
		Run struct {
			//Text string `xml:",chardata"`
			ID string `xml:"id,attr"`
		} `xml:"Run"`
		Platform struct {
			//Text         string `xml:",chardata"`
			UsesEmulator string `xml:"usesEmulator,attr"`
		} `xml:"Platform"`
		Region    string `xml:"Region"`
		Variables string `xml:"Variables"`
	} `xml:"Metadata"`
	Offset         string `xml:"Offset"`
	AttemptCount   string `xml:"AttemptCount"`
	AttemptHistory struct {
		//Text    string `xml:",chardata"`
		Attempt []struct {
			//Text            string `xml:",chardata"`
			ID              string `xml:"id,attr"`
			Started         string `xml:"started,attr"`
			IsStartedSynced string `xml:"isStartedSynced,attr"`
			Ended           string `xml:"ended,attr"`
			IsEndedSynced   string `xml:"isEndedSynced,attr"`
			RealTime        string `xml:"RealTime"`
			GameTime        string `xml:"GameTime"`
		} `xml:"Attempt"`
	} `xml:"AttemptHistory"`
	Segments struct {
		//Text    string `xml:",chardata"`
		Segment []struct {
			//Text       string `xml:",chardata"`
			Name       string `xml:"Name"`
			Icon       string `xml:"Icon"`
			SplitTimes struct {
				//Text      string `xml:",chardata"`
				SplitTime []struct {
					//Text     string `xml:",chardata"`
					Name     string `xml:"name,attr"`
					RealTime string `xml:"RealTime"`
					GameTime string `xml:"GameTime"`
				} `xml:"SplitTime"`
			} `xml:"SplitTimes"`
			BestSegmentTime struct {
				//Text     string `xml:",chardata"`
				RealTime string `xml:"RealTime"`
				GameTime string `xml:"GameTime"`
			} `xml:"BestSegmentTime"`
			SegmentHistory struct {
				//Text string `xml:",chardata"`
				Time []struct {
					//Text     string `xml:",chardata"`
					ID       string `xml:"id,attr"`
					RealTime string `xml:"RealTime"`
					GameTime string `xml:"GameTime"`
				} `xml:"Time"`
			} `xml:"SegmentHistory"`
		} `xml:"Segment"`
	} `xml:"Segments"`
	AutoSplitterSettings struct {
		//Text           string `xml:",chardata"`
		Version        string `xml:"Version"`
		ScriptPath     string `xml:"ScriptPath"`
		CustomSettings string `xml:"CustomSettings"`
	} `xml:"AutoSplitterSettings"`
}

func (run *XMLRun) getSegmentNames() []string {
	var names []string
	for _, segment := range run.Segments.Segment {
		names = append(names, segment.Name)
	}
	return names
}

func (run *XMLRun) getPBTimes() []string {
	var times []string
	for _, segment := range run.Segments.Segment {
		for _, splittime := range segment.SplitTimes.SplitTime {
			if splittime.Name == "Personal Best" {
				times = append(times, splittime.RealTime)
				//else the config says gametime
			}
		}
	}
	return times
}

/*
	TODO:
		check game names and combine categories into one section
*/

func main() {
	/*
		loop through splits folder
		read each file
		write markdown for each file
	*/
	dirent, err := os.ReadDir("splits")
	if err != nil {
		log.Fatal(err)
	}

	readMe, err := os.OpenFile("readme.md", os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dirent {
		runData := &XMLRun{}

		fileData, err := os.ReadFile("splits/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		xml.Unmarshal(fileData, &runData)

		readMe.WriteString("## " + runData.GameName + "\n")
		readMe.WriteString("### " + runData.CategoryName + " Attempts: " + runData.AttemptCount + "\n")
		names := runData.getSegmentNames()
		times := runData.getPBTimes()

		for i, n := range names {
			readMe.WriteString("- " + n + " - " + times[i] + "\n")
		}
	}

}
