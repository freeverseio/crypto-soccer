using System.Numerics;
public class Serialization {  
    public double Add(double num1, double num2) {  
        return num1 + num2;  
    }  
    public BigInteger AddBN(BigInteger x, BigInteger y) {  
        return BigInteger.Add(x, y);  
    }

    public uint getCurrentShirtNum(BigInteger playerState) {
        return (uint) ((playerState >> 43) & 31);
    }

    // public uint getSkill(BigInteger encodedSkills, skillIdx uint8){
    //     return and(right(encodedSkills, uint(skillIdx)*20), 1048575) // 1048575 = 2**20 - 1
    // }

    // public uint getBirthDay(BigInteger encodedSkills) {
    //     return and(right(encodedSkills, 100), 65535)
    // }
}  