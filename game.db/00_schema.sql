CREATE TABLE player_props (
    player_id TEXT NOT NULL,
    player_name TEXT NOT NULL,
    PRIMARY KEY(player_id)
);

CREATE TABLE team_props (
    team_id TEXT NOT NULL,
    team_name TEXT NOT NULL,
    team_manager_name TEXT NOT NULL,
    PRIMARY KEY(team_id)
);

CREATE OR REPLACE FUNCTION public.playernamebyplayerid(playerid text)
 RETURNS text
 LANGUAGE plpgsql
 STABLE
AS $function$
declare
   result_name text;
begin
   select player_name 
   into result_name
   from player_props
   where player_id = playerid;
   
   return result_name;
end;
$function$
;

CREATE OR REPLACE FUNCTION public.teamnamebyteamid(teamid text)
 RETURNS text
 LANGUAGE plpgsql
 STABLE
AS $function$
declare
   result_name text;
begin
   select team_name 
   into result_name
   from team_props
   where team_id = teamid;
   
   return result_name;
end;
$function$
;

CREATE OR REPLACE FUNCTION public.teammanagernamebyteamid(teamid text)
 RETURNS text
 LANGUAGE plpgsql
 STABLE
AS $function$
declare
   result_name text;
begin
   select team_manager_name 
   into result_name
   from team_props
   where team_id = teamid;
   
   return result_name;
end;
$function$
;