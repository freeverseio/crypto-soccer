using Microsoft.VisualStudio.TestTools.UnitTesting;

namespace NewMSTestProject
{
    [TestClass]
    public class UnitTest1
    {
        [TestMethod]
        public void Test_AddMethod() {  
            BasicMaths bm = new BasicMaths();  
            double res = bm.Add(10, 10);  
            Assert.AreEqual(res, 20);  
        }  
    }
}
