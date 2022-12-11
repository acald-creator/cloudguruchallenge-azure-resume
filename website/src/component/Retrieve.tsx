import React from "react"
import useSWR from "swr"

const fetcher = (url: string) => fetch(url).then((res) => res.json());

function Retrieve() {
    const { data, error, isLoading } = useSWR(
        "https://resume.acaldwell.dev/resume.json",
        fetcher
    );

    if (error) return <div>failed to load</div>
    if (isLoading) return <div>loading...</div>

    return (
        <div className="prose container mx-auto">
            <h1 className="text-3xl font-bold underline">{data.basics.name}</h1>
            <p>{data.basics.label}</p>
            <p className="font-bolder italic">You are visitor: *placeholder* </p>
        </div>
    );
}

export default function Display() {
    return (
        <Retrieve />
    )
}