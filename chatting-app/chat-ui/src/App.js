import React, { useEffect } from "react";
import { Route, Routes, useNavigate } from "react-router-dom";
import { setNavigate } from "./services/navigation";
import Login from "./services/Login";
import MainChannel from "./services/MainChannel";

const App = () => {
	const navigate = useNavigate();

	useEffect(() => {
		setNavigate(navigate); // Set the global navigate function
	  }, [navigate]);

	return (
			<Routes>
				<Route path="/channels" element={<MainChannel />} />
				<Route path="/channels/:channelId" element={<MainChannel />} />
				<Route path="/" element={<Login />} />
        		<Route path="/login" element={<Login />} />
			</Routes>
	);
};

export default App;