const APP_VERSION: &str = env!("CARGO_PKG_VERSION");

fn main() {
    println!("Hello, rs! [Version: {}]", APP_VERSION);
}

pub struct User {
    id: usize,
}

impl User {
    pub fn new() -> User {
        User { id: 0 }
    }

    pub fn print(&self) {
        println!("User: id = {}", self.id)
    }
}

impl Default for User {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    #[test]
    fn foo() {
        assert_eq!(2 * 3, 5)
    }
}
