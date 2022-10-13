fn main() {
    let mut page = 0;
    let agent = ureq::AgentBuilder::new().build();
    let mut users: Vec<User> = Vec::with_capacity(10_000);
    let file = std::fs::File::create("results.json").expect("Failed to open json file!");
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
        for user in data.players {
            users.push(User {
                id: user.id,
                msgs: user.message_count,
                xp: user.xp,
                level: user.level,
            });
        }
        if page >= 100 {
            break;
        }
        println!("Fetched page {} out of 100", page);
        page += 1;
        std::thread::sleep(std::time::Duration::from_secs_f32(1.2))
    }
    serde_json::to_writer_pretty(file, &users).expect("Failed to serialize users as json!");
}

#[derive(serde::Serialize)]
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
