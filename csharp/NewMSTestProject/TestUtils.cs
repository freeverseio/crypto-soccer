using System.Numerics;
using System;
using System.IO;
using Newtonsoft.Json;

public class TestUtils {  
    public dynamic LoadJson(string filename)
    {
        string workingDirectory = Environment.CurrentDirectory;
        string codeDirectory = Directory.GetParent(workingDirectory).Parent.Parent.FullName;
        string[] paths = {codeDirectory, "testdata", filename};
        string jsonFile = Path.Combine(paths);

        Console.WriteLine(jsonFile);

        dynamic array;
        using (StreamReader r = new StreamReader(jsonFile))
        {
            string json = r.ReadToEnd();
            array = JsonConvert.DeserializeObject(json);
            Console.WriteLine("Read a total of {0} tests", array.Count);
        }
        return array;
    }
}  