import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { fetchWithToken } from "./api";

const ChannelsList = ({ selectedChannel, setSelectedChannel }) => {

	const { channelId } = useParams();
	const [channels, setChannels] = useState([]);
	const [newChannelName, setNewChannelName] = useState("");
	const [menuPosition, setMenuPosition] = useState(null);
	const [selectedChannelForAction, setSelectedChannelForAction] = useState(null);

	useEffect(() => {
		if (channelId) {
			console.log(channelId)
			const channel = channels.find(
				(channel) => channel.ID === parseInt(channelId),
			);
			if (channel) {
				setSelectedChannel({ ChannelName: channel.ChannelName, ID: parseInt(channelId) });
			}
		}
	}, [channelId, channels]);

	useEffect(() => {
		const fetchChannels = async () => {
			const response = await fetchWithToken("/channel/");
			if (!response) {
				console.error("Failed to fetch channels due to unauthorized or network error.");
				return; // Exit early if response is null
			  }
		  
			  try {
				const data = await response.json();
				setChannels(data || []);
			  } catch (error) {
				console.error("Error parsing JSON:", error);
			  }
		};
		fetchChannels();
	}, []);

	const handleAddChannel = async () => {
		const response = await fetchWithToken("/channel/", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ channelName: newChannelName }),
		});
		console.log(response)

		if (response.ok) {
			const newChannel = await response.json();
			setChannels([...channels, { ID: newChannel.ID, ChannelName: newChannelName }]);
			setNewChannelName("");
		}
	};

	const handleUpdateChannel = async () => {
		const newName = prompt("Enter new channel name", selectedChannelForAction.ChannelName);
		if (newName) {
		  	const response = await fetchWithToken(`/channel/${selectedChannelForAction.ID}`, {
				method: "PUT",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ channelName: newName }),
			});
	
			if (response.ok) {
				setChannels((prevChannels) =>
				prevChannels.map((channel) =>
					channel.ID === selectedChannelForAction.ID
					? { ...channel, ChannelName: newName }
					: channel
				)
				);
				setMenuPosition(null);
			}
		}
	};
	
	const handleDeleteChannel = async () => {
		if (window.confirm("Are you sure you want to delete this channel?")) {
			const response = await fetchWithToken(`/channel/${selectedChannelForAction.ID}`, {
				method: "DELETE",
			});
		
			if (response.ok) {
				setChannels((prevChannels) =>
				prevChannels.filter((channel) => channel.ID !== selectedChannelForAction.ID)
				);
				setMenuPosition(null);
			}
		}
	};

	const handleContextMenu = (event, channel) => {
		event.preventDefault();
		setSelectedChannelForAction(channel);
		setMenuPosition({ x: event.clientX, y: event.clientY })
	}

	const handleCloseMenu = () => {
		setMenuPosition(null)
	}

	return (
		<div className="flex flex-col h-full bg-gray-100 border-r">
			<div className="bg-gray-700 text-white p-2">Channels</div>
			<div className="overflow-y-auto flex-grow p-4">
				{channels ? (
					<ul className="w-full">
						{channels.map((channel) => (
							<li
								key={channel.ID}
								onClick={() => setSelectedChannel(channel)}
								onContextMenu={(event) => handleContextMenu(event, channel)}
								className={`p-2 rounded-md w-full cursor-pointer ${
									parseInt(channelId) === channel.ID ? "bg-blue-500 text-white" : "hover:bg-gray-200"
								}`}
							>
								{channel.ChannelName}
							</li>
						))}
					</ul>
				) : (
					<div className="text-center text-gray-600">Please add a Channel</div>
				)}
			</div>
			<div className="flex flex-col p-4">
				<input
					type="text"
					value={newChannelName}
					onChange={(e) => setNewChannelName(e.target.value)}
					placeholder="New channel..."
					className="mb-4 p-2 w-full border rounded-md bg-white"
				/>
				<button
					onClick={handleAddChannel}
					className="p-2 bg-blue-500 text-white rounded-md"
				>
					Add Channel
				</button>
			</div>
		</div>
	);
};

export default ChannelsList;