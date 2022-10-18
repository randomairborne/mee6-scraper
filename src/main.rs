use std::io::Write;


fn main() {
    let agent = ureq::AgentBuilder::new().build();
    let mut page = 0;
    let mut users: Vec<User> = Vec::new();
    let file = std::fs::OpenOptions::new()
        .create(true)
        .write(true)
        .open("results.json")
        .expect("Failed to open results.json!");
    loop {
        let req = match agent
            .get(&format!(
                "https://mee6.xyz/api/plugins/levels/leaderboard/{}",
                std::env::args()
                    .nth(1)
                    .unwrap_or("302094807046684672".to_string())
            ))
            .query("page", &page.to_string())
            .call()
        {
            Ok(val) => val,
            Err(e) => {
                println!("Error fetching: {}", e);
                break;
            }
        };
        let data = match req.into_json::<Root>() {
            Ok(val) => val,
            Err(e) => {
                println!("Error deserializing: {e:?}");
                break;
            }
        };
        let last = data.players[data.players.len() - 1].clone();
        for user in data.players {
            users.push(User {
                id: user.id,
                msgs: user.message_count,
                xp: user.xp,
                level: user.level,
            });
        }
        if last.level < 5 {
            break;
        }
        page += 1;
        print!("\rCurrent user level: {} ({} total users) ", last.level, users.len());
        std::io::stdout().flush().ok();
        std::thread::sleep(std::time::Duration::from_secs(1));
    }
    serde_json::to_writer_pretty(file, &users).expect("Failed to serialize users as json!");
}

#[derive(Default, Debug, Clone, PartialEq, serde::Deserialize, serde::Serialize)]
struct User {
    id: String,
    msgs: i64,
    xp: i64,
    level: i64,
}

#[derive(Default, Debug, Clone, PartialEq, serde::Deserialize)]
pub struct Root {
    pub page: i64,
    pub players: Vec<Player>,
}

#[derive(Default, Debug, Clone, PartialEq, serde::Deserialize)]
pub struct Player {
    pub id: String,
    pub level: i64,
    pub message_count: i64,
    pub xp: i64,
}
