Assets 
    is EncodingSkills, EncodingState, EncodingIDs 
    
Market 
    points to Assets
    
Engine 
    is EngineLib, EncodingMatchLogPart3
    points to EnginePreComp

EngineLib 
    is EncodingSkills    
    
EnginePreComp 
    is EngineLib, EncodingMatchLogPart1, SortValues 

Evolution 
    is EncodingMatchLog, EncodingSkills, EngineLib, EncodingTPAssignment, EncodingSkillsSetters {
    
Championships 
    is SortIdxs, EncodingSkills
    points to Engine
    
Friendlies 
    is SortIdxsAnySize 

Updates
    points to Assets

EncodingMatchLog 
    is EncodingMatchLogPart1, EncodingMatchLogPart2, EncodingMatchLogPart3    