package ohbot

const (
	settingsDef = `<SettingList>
    <Setting Name="DefaultEyeShape" Value="eyeball" /> 
    <Setting Name="DefaultSpeechSynth" Value="" /> 
    <Setting Name="DefaultVoice" Value="" /> 
    <Setting Name="DefaultLang" Value="en-gb" /> 
    <Setting Name="SpeechDBFile" Value="ohbotData/OhbotSpeech.csv" /> 
    <Setting Name="EyeShapeList" Value="ohbotData/ohbot.obe" /> 
    <Setting Name="MotorDefFile" Value="ohbotData/MotorDefinitionsv21.omd" />
</SettingList>`

	motorDef = `<Motors>
  <Motor Name="HeadTurn" Min="0" Max="1000" Motor="1" Speed="40" Reverse="False" Acceleration="60" RestPosition="5" Avoid="" />
  <Motor Name="HeadNod" Min="140" Max="700" Motor="0" Speed="0" Reverse="True" Acceleration="60" RestPosition="5" Avoid="" />
  <Motor Name="EyeTurn" Min="380" Max="780" Motor="2" Speed="0" Reverse="False" Acceleration="0" RestPosition="5" Avoid="" />
  <Motor Name="EyeTilt" Min="520" Max="920" Motor="6" Speed="0" Reverse="False" Acceleration="30" RestPosition="5" Avoid="" />
  <Motor Name="TopLip" Min="0" Max="550" Motor="4" Speed="0" Reverse="True" Acceleration="0" RestPosition="5" Avoid="BottomLip" />
  <Motor Name="BottomLip" Min="0" Max="550" Motor="5" Speed="0" Reverse="True" Acceleration="0" RestPosition="5" Avoid="TopLip" />
  <Motor Name="LidBlink" Min="0" Max="350" Motor="3" Speed="0" Reverse="False" Acceleration="0" RestPosition="10" Avoid="" />
  <Motor Name="MouthOpen" Min="80" Max="460" Motor="7" Speed="0" Reverse="False" Acceleration="0" RestPosition="10" Avoid="" />
</Motors>`

	speechDef = `set,variable,phrase
1,1,what is your command?
1,2,how can I help?
1,3,what's up?
1,4,ready for input
2,1,Let me think
2,2,Just a second
2,3,give me a moment
2,4,let me get you an answer
3,1,Hello my name is Ohbot`
)
