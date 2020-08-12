using System.Numerics;
using System;
using System.IO;
using Newtonsoft.Json;

public class TestUtils {  

    public class Item
    {
        public string encodedSkills;
        public int[] skills;
        public int birthday;
        public bool isSpecial;
        public string playerIdFromSkills;
        public string internalPlayerId;
        public int potential;
        public int forwardness;
        public int leftishness;
        public int aggressiveness;
        public bool alignedEndOfFirstHalf;
        public bool redCardLastGame;
        public int gamesNonStopping;
        public int injuryWeeksLeft;
        public bool substitutedFirstHalf;
        public int sumOfSkills;
        public int generation;
        public bool outOfGameFirstHalf;
        public bool yellowCardFistHalf;
    }

    public void LoadJson()
    {
        string workingDirectory = Environment.CurrentDirectory;
        string codeDirectory = Directory.GetParent(workingDirectory).Parent.Parent.FullName;
        string[] paths = {codeDirectory, "testdata", "encodingSkillsTestData.json"};
        string jsonFile = Path.Combine(paths);

        Console.WriteLine(jsonFile);

        using (StreamReader r = new StreamReader(jsonFile)) {}
        // {
        //     string json = r.ReadToEnd();
        //     // List<Item> items = JsonConvert.DeserializeObject<List<Item>>(json);
        //     // dynamic array = JsonConvert.DeserializeObject(json);
        //     // Console.WriteLine(json);
        //     // JsonTextReader reader = new JsonTextReader(new StringReader(json));
        //     // while (reader.Read())
        //     // {
        //     //     Console.WriteLine(reader);

        //     //     // if (reader.Value != null)
        //     //     // {
        //     //     //     Console.WriteLine("Token: {0}, Value: {1}", reader.generation, reader.generation);
        //     //     // }
        //     // }        
        // }
    }

}  