using Microsoft.VisualStudio.TestTools.UnitTesting;
using System.Numerics;
using System;

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

        [TestMethod]
        public void getCurrentShirtNum() {  
            Serialization serial = new Serialization();
            uint shirt = 13;
            BigInteger state = new BigInteger(shirt * Math.Pow(2,43));
            uint res = serial.getCurrentShirtNum(state);  
            Assert.AreEqual(res, (uint) shirt);  
        }  
    }
}
