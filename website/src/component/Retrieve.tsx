import React, { Suspense } from "react"
import useSWR from "swr"

const fetcher = (url: string) => fetch(url).then((res) => res.json());

function Retrieve() {
    const { data, error, isLoading } = useSWR(
        "https://resume.acaldwell.dev/resume.json",
        fetcher,
        { suspense: true }
    );

    if (error) return <div>Oops! Failed to load.</div>
    if (isLoading) return <div>Loading...</div>

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
        <Suspense fallback={<h1>Loading page...</h1>}>
            <Retrieve />
        </Suspense>
    )
}