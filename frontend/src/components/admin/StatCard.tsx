interface StatCard {
    title: string;
    value: number;
    icon: React.ReactNode;
    bgColor: string;
    iconColor: string;
}

interface StatsCardsProps {
    totalStages: number;
    totalCapacity: number;
    occupiedPlaces: number;
    availablePlaces: number;
}

export default function StatsCards({
                                       totalStages,
                                       totalCapacity,
                                       occupiedPlaces,
                                       availablePlaces
                                   }: StatsCardsProps) {
    const stats: StatCard[] = [
        {
            title: "Total stages",
            value: totalStages,
            bgColor: "bg-blue-100",
            iconColor: "text-blue-600",
            icon: (
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                </svg>
            )
        },
        {
            title: "Capacité totale",
            value: totalCapacity,
            bgColor: "bg-green-100",
            iconColor: "text-green-600",
            icon: (
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
            )
        },
        {
            title: "Places occupées",
            value: occupiedPlaces,
            bgColor: "bg-orange-100",
            iconColor: "text-orange-600",
            icon: (
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
                </svg>
            )
        },
        {
            title: "Places disponibles",
            value: availablePlaces,
            bgColor: "bg-purple-100",
            iconColor: "text-purple-600",
            icon: (
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192L5.636 18.364M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-5 0a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
            )
        }
    ];

    return (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            {stats.map((stat, index) => (
                <div key={index} className="bg-white p-6 rounded-lg shadow-sm border">
                    <div className="flex items-center">
                        <div className={`p-2 ${stat.bgColor} rounded-lg`}>
                            <div className={stat.iconColor}>
                                {stat.icon}
                            </div>
                        </div>
                        <div className="ml-4">
                            <p className="text-sm font-medium text-gray-500">{stat.title}</p>
                            <p className="text-2xl font-semibold text-gray-900">{stat.value}</p>
                        </div>
                    </div>
                </div>
            ))}
        </div>
    );
}
