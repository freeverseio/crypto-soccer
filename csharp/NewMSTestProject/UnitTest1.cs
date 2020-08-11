using Microsoft.VisualStudio.TestTools.UnitTesting;
using System.Numerics;

namespace NewMSTestProject
{
    [TestClass]
    public class UnitTest1
    {
        [TestMethod]
        public void Test_AddMethod() {  
            Serialization serial = new Serialization();  
            double res = serial.Add(10, 20);  
            Assert.AreEqual(res, 30);  
        }  
        [TestMethod]
        public void Test_AddBNMethod() {  
            Serialization serial = new Serialization();  
            BigInteger res = serial.AddBN(new BigInteger(10), new BigInteger(20));  
            Assert.AreEqual(res, new BigInteger(30));  
        }  
    }
}
