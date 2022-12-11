import React from 'react';
import logo from './logo.svg';
import './App.css';
import useSWR from "swr"

const fetcher = (url) => fetch(url).then((res) => res.json());

export default function App() {
  const { data, error } = useSWR(
    "https://resume.acaldwell.dev/resume.json",
    fetcher
  );

  if (error) return "An error has occurred.";
  if (!data) return "Loading...";

  return (
    <React.Fragment>
      <div>
        <h1>{data.basics.name}</h1>
        <p>{data.basics.label}</p>
        <p>You are visitor: </p>
      </div>
    </React.Fragment>
  );
}
