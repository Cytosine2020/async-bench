use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::{TcpListener, TcpStream};

use std::error::Error;

const CLIENT_COUNT : usize = 1024;
const MSG_COUNT : usize = 1024;
const BUFFER_SIZE: usize = 16;

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let mut handle = vec![];

    let mut port = 54321;

    let listener = loop {
        match TcpListener::bind(format!("127.0.0.1:{}", port)).await {
            Ok(listener) => break Ok(listener),
            Err(e) => 
                if port < 54421 {
                    break Err(e);
                },
        }

        port += 1;
    }?;

    for _ in 0..CLIENT_COUNT {
        handle.push(tokio::spawn(async move {
            async fn client(port :u16) -> Result<(), Box<dyn Error>> {
                let mut socket = TcpStream::connect(format!("127.0.0.1:{}", port)).await?;

                let mut send = vec![0u8; BUFFER_SIZE];
                let mut recv = vec![0u8; BUFFER_SIZE];

                for _ in 0..MSG_COUNT {
                    send.fill_with(rand::random);

                    socket.write_all(&send).await?;

                    let n = socket.read(&mut recv).await?;
                    assert_eq!(n, BUFFER_SIZE);
                    assert_eq!(send, recv);
                }

                Ok(())
            }

            client(port).await.unwrap();
        }));
    }

    for _ in 0..CLIENT_COUNT {
        let (socket, _) = listener.accept().await?;

        handle.push(tokio::spawn(async move {
            async fn server(mut socket: TcpStream) -> Result<(), Box<dyn Error>> {
                let mut buf = vec![0; BUFFER_SIZE];

                loop {
                    let n = socket.read(&mut buf).await?;

                    if n == 0 {
                        return Ok(());
                    }

                    socket.write_all(&buf[0..n]).await.expect("failed to write data to socket");
                }
            }

            server(socket).await.unwrap();
        }));
    };

    futures::future::join_all(handle).await;

    Ok(())
}
