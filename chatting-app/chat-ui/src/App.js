import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Login from "./Login";
import MainChat from "./MainChat";

const App = () => {
	return (
		<Router>
			<Routes>
				<Route path="/chat" element={<MainChat />} />
				<Route path="/chat/:channelId" element={<MainChat />} />
				<Route path="/" element={<Login />} />
        <Route path="/login" element={<Login />} />
			</Routes>
		</Router>
	);
};

export default App;