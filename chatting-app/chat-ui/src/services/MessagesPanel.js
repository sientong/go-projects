import React, { useState, useEffect, useRef } from "react";
// import MessageEntry from "./MessageEntry";
// import { fetchWithToken } from "./api";

const MessagesPanel = ({ selectedChannel }) => {
    const [messages, setMessages] = useState([]);
    const websocketRef = useRef(null);

    const printMessage = (message) => {
        const container = document.querySelector(".messages-container");
        const el = document.createElement("p");
        el.innerHTML = message;
        container.append(el);
    };

    const sendMessage = () => {
        const input = document.querySelector(".input-message");
        const messageRaw = input.value;

        if (websocketRef.current && messageRaw.trim() !== "") {
            console.log(websocketRef.current)
            websocketRef.current.send(JSON.stringify({
                Message: messageRaw
            }));


            const formattedMessage = `<b>me</b>: ${messageRaw}`;
            printMessage(formattedMessage);

            input.value = "";
        }
    };

    useEffect(() => {
        if (!selectedChannel) return;

        const connectWebSocket = () => {
            const username = localStorage.getItem("username")
            websocketRef.current = new WebSocket(`ws://localhost:8080/ws?username=${username}&channelID=${selectedChannel.ID}`);

            websocketRef.current.onopen = () => {
                console.log("WebSocket connected");
            };

            websocketRef.current.onmessage = (event) => {
                const newMessage = JSON.parse(event.data);
                console.log(event.data)
                setMessages((prevMessages) => [...prevMessages, newMessage]);

                var formattedMessage;
                if (newMessage.Type === "Leave") {
                    formattedMessage = `<b>${newMessage.From}</b> leaves the channel`;
                } else if (newMessage.Type === "New User"){
                    formattedMessage = `<b>${newMessage.From}</b> joins the channel`;
                } else {
                    formattedMessage = `<b>${newMessage.From}</b>: ${newMessage.Message}`;
                } 
                printMessage(formattedMessage);
                
            };

            websocketRef.current.onclose = () => {
                console.log("WebSocket closed, reconnecting...");
                setTimeout(connectWebSocket, 5000); // Reconnect after 5 seconds
            };

            websocketRef.current.onerror = (error) => {
                console.error("WebSocket error:", error);
                websocketRef.current.close();
            };
        };

        connectWebSocket();

        return () => {
            if (websocketRef.current) {
                websocketRef.current.close();
            }
        };
    }, [selectedChannel]);

    return (
        <div className="flex flex-col h-full">
            {selectedChannel && (
                <div className="bg-gray-700 text-white p-2">
                    Messages for {selectedChannel.name}
                </div>
            )}
            <div
                className={`messages-container overflow-auto flex-grow ${
                    selectedChannel && messages.length === 0
                        ? "flex items-center justify-center"
                        : ""
                }`}
            >
                {selectedChannel ? (
                    messages.length > 0 ? (
                        messages.map((message) => (
                            <div className="p-2 border-b">
                            </div>
                        ))
                    ) : (
                        <div className="text-center text-gray-600">
                            No messages yet! Why not send one?
                        </div>
                    )
                ) : (
                    <div className="p-2">Please select a channel</div>
                )}
            </div>
            {selectedChannel && (
                <div className="p-2">
                    <input type="text" className="input-message border p-2 w-full" placeholder="Type your message..." />
                    <button
                        className="bg-blue-500 text-white p-2 mt-2"
                        onClick={sendMessage}
                    >
                        Send Message
                    </button>
                </div>
            )}
        </div>
    );
};

export default MessagesPanel;
