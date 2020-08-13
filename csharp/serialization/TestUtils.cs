using Microsoft.VisualStudio.TestTools.UnitTesting;
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

    public void AssertDecodePlayerStateOK(dynamic test)
    {
        Serialization serial = new Serialization();
        BigInteger encoded;
        bool succeeded = BigInteger.TryParse((string) test.encodedState, out encoded);
        Assert.AreEqual(true, succeeded);
        Assert.AreEqual((ulong) test.currentTeamId, serial.getCurrentTeamId(encoded));  
        Assert.AreEqual((uint) test.currentShirtNum, serial.getCurrentShirtNum(encoded));  
        Assert.AreEqual((ulong) test.prevPlayerTeamId, serial.getPrevPlayerTeamId(encoded));  
        Assert.AreEqual((ulong) test.lastSaleBlocknum, serial.getLastSaleBlock(encoded));  
        Assert.AreEqual((bool) test.isInTransit, serial.getIsInTransit(encoded));  
    }

}  