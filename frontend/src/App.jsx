import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Signup from './Signup.jsx';
import Signin from './Signin.jsx';
import UrlTable  from "./UrlTable";
import HomePage from './HomePage.jsx';
import BeanUrlAppBar from "./BeanAppbar.jsx";
import InputBar from './InputBar.jsx';
import CollapsibleTable from './AdminPanel.jsx';
import axios from 'axios';

function App() {
    const [authenticatedUser, setAuthenticatedUser] = useState(false);
    const [authenticatedAdmin, setAuthenticatedAdmin] = useState(false);
    const [urlDetailRows, setUrlDetailRows] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await axios.get("http://localhost:3000/user/me", {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("userToken")
                    }
                });

                const userData = response.data.data.user;
                if (userData && userData.userRole === "user") {
                    setAuthenticatedUser(true);
                } else if (userData && userData.userRole === "admin") {
                    setAuthenticatedAdmin(true);
                } else {
                    setAuthenticatedUser(false);
                }
            } catch (error) {
                console.error("Error while reading claims from jwt", error);
                setAuthenticatedUser(false);
                setAuthenticatedAdmin(false);
            }
        };

        fetchData();
    }, []);

    useEffect(() => {
        const fetchUrlData = async () => {
            try {
                const response = await fetch("http://localhost:3000/url-shortener/resolutions/analytics", {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": "Bearer " + localStorage.getItem("userToken")
                    }
                });

                const data = await response.json();
                console.log("API Response:", data);
                setUrlDetailRows(data);
            } catch (error) {
                console.error("Error fetching data:", error);
            }
        };

        fetchUrlData();
    }, []);


    return (
        <div className="app-container"> {/* Wrap the entire app */}
            <Router>
                <BeanUrlAppBar/>
                <Routes>
                    <Route path={"/signup"} element={<Signup />} />

                    {/* Only allow access to the signin route if not authenticatedUser */}
                    {!authenticatedUser && <Route path={"/signin"} element={<Signin />} />}
                    {/* If Authenticated user tries to go to signin route manually.. route him/her to HomePage */}
                    {authenticatedUser && <Route path={"/signin"} element={<HomePage urlDetailRows={urlDetailRows} setUrlDetailRows={setUrlDetailRows} />} />}
                    
                    {/* Secure routes for user*/}
                    {authenticatedUser ? (
                        <>
                        <Route path={"/user"} element={
                            <div>
                                <InputBar urlDetailRows={urlDetailRows} setUrlDetailRows={setUrlDetailRows} />
                                <UrlTable urlDetailRows={urlDetailRows} setUrlDetailRows={setUrlDetailRows} />
                            </div>
                        } />
                        </>
                    ) : null}


                    {/* Secure routes for admin*/}
                    {authenticatedAdmin ? (
                        <>
                        <Route path={"/admin"} element={<CollapsibleTable/>}/>
                        </>
                    ) : null}


                    {/* Redirect to home if trying to access secure routes while not authenticatedUser */}
                    {!authenticatedUser && <Route path={"/*"} element={<HomePage urlDetailRows={urlDetailRows} setUrlDetailRows={setUrlDetailRows} />} />}
                    
                </Routes>
            </Router>
        </div>
    );
}

export default App;
