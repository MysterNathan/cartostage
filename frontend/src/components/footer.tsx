"use client";

export default function Footer() {
    return (
    <footer className="bg-white border-t">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
            <div className="text-center text-gray-600 text-sm">
                <p>© {new Date().getFullYear()} Carte des stages - Tous droits réservés</p>
            </div>
        </div>
    </footer>
    )
}