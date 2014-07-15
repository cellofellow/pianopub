# PianoPub TODO

Terminology:
* CONNECTED USERS means individual connections with unique login IDs.
* VOTING USERS means CONNECTED USERS who have voted on a proposal. All cases
  where decisions are maded by VOTING USERS is must have a time limit defined.

List of features.

## Admin account that can:

* Set up Pandora credentials.
* Set up MPD details.
* Manage user accounts (add/remove privileges, ban, etc).

## Basic usage.

Two modes: filesystem mode and Pandora mode. Possibly in the future add
ShoutCAST and Xiph.org and rdio too.

### Voting

Any user can at any time propose a voting action. Voting actions must then be
approved by majority or supermajority, depending on the action, of either
VOTING USERS or CONNECTED USERS.

All proposals will send a notification of the vote to all connections. Votes by
definition only allow binary responses. Votes with a CONNECTED USERS
requirement are assumed to be in the negative when no vote is received.

A user earns one rep point for each proposal they make that passes.

### Filesystem mode.

Default startup mode. Frontends MPDs music library. Initial playback mode is
full shuffle of entire MPD library.

User can propose PAUSE PLAYBACK that requires a majority of VOTING USERS within
1 minute.

User can propose SKIP SONG that requires a majority of VOTING USERS within 30
seconds to pass.

User can propose DELETE SONG that deletes the currently playing song. Requires
a supermajority of CONNECTED USERS.

User can create a NEW PLAYLIST by submitting a list of song IDs.

User can MODIFY their own playlist at will.

User can FORK anothers playlist at will.

User can submit a MERGE proposal to another playlist requiring users approval
(git-style merge requests).

User can propose a PLAY PLAYLIST. A majority of VOTING USERS within 1 MINUTE is
required to switch to the playlist.

### Pandora Mode

Mode that connects to Pandora.com and streams algorithmically-generated radio
stations. Uses a single Pandora account.

Switching modes starts with a user proposal followed by a majority of CONNECTED
USERS vote.

User can propose PAUSE PLAYBACK that requires a majority of VOTING USERS within
1 minute.

Starting condition for Pandora mode is to play the special "Quickmix" station.

User can propose SKIP SONG that requires a majority of VOTING USERS within 30
seconds.

User can propose LIKE SONG that marks the song liked with Pandora. Requires a
majority of CONNECTED USERS. Connected users is used instead of voting users
because this can sharply skew a stations algorithm.

User can propose DISLIKE SONG that marks a song banned with Pandora. Requires a
majority of CONNECTED USERS, for similar reasons as above.

User can propose TIRED OF SONG, which marks song as temporarily banned with
Pandora. Requires a majority of VOTING USERS within 30 seconds.

User can propose SWITCH STATION which requires a majority of VOTING USERS
within 1 minute.

User can search for artists, albums, and songs on Pandora that can be used as
station seeds.

User can propose CREATE STATION with a seed from the search. Requires a
majority of CONNECTED USERS.

User can propose STATION QUICK MIX which will set specified station to the
quick mix. Also the opposite. Requires a majority of CONNECTED USERS.

User can propose MODIFY STATION to add a seed to a station. Requires a majority
of CONNECTED USERS.

User can propose DELETE STATION. Requires a majority of CONNECTED USERS.


