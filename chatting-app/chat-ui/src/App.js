import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Login from "./services/Login";
import MainChannel from "./services/MainChannel";

const App = () => {
	return (
		<Router>
			<Routes>
				<Route path="/channels" element={<MainChannel />} />
				<Route path="/channels/:channelId" element={<MainChannel />} />
				<Route path="/" element={<Login />} />
        <Route path="/login" element={<Login />} />
			</Routes>
		</Router>
	);
};

export default App;