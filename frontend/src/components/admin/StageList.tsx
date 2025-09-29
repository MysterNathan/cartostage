// components/admin/StagesList.tsx
import React from 'react';
import { Stage } from '@/types/stage';

interface StagesListProps {
    stages: Stage[];
    onEdit: (stage: Stage) => void;
    loading?: boolean;
}

const StagesList: React.FC<StagesListProps> = ({
                                                   stages,
                                                   onEdit,
                                                   loading = false
                                               }) => {
    const getFiliereStyles = (filiere: string) => {
        const styles = {
            CCST: 'bg-blue-100 text-blue-800',
            SN: 'bg-green-100 text-green-800',
            MELEC: 'bg-yellow-100 text-yellow-800',
            TMA: 'bg-red-100 text-red-800'
        };
        return styles[filiere as keyof typeof styles] || 'bg-purple-100 text-purple-800';
    };

    const getParcoursStyles = (parcours: string) => {
        const styles = {
            scolaire: 'bg-blue-50 text-blue-700 border-blue-200',
            apprentissage: 'bg-orange-50 text-orange-700 border-orange-200',
            mixte: 'bg-purple-50 text-purple-700 border-purple-200'
        };
        return styles[parcours as keyof typeof styles] || 'bg-gray-50 text-gray-700 border-gray-200';
    };

    const getNiveauStyles = (niveau: string) => {
        const styles = {
            '2de': 'bg-green-100 text-green-800',
            '1re': 'bg-yellow-100 text-yellow-800',
            'Tle': 'bg-red-100 text-red-800'
        };
        return styles[niveau as keyof typeof styles] || 'bg-gray-100 text-gray-800';
    };

    const getCapacityStyles = (filled: number, total: number, disponibles: number) => {
        if (filled >= total) return { text: 'text-red-600', bar: 'bg-red-500' };
        if (disponibles <= 1) return { text: 'text-orange-600', bar: 'bg-orange-500' };
        return { text: 'text-green-600', bar: 'bg-green-500' };
    };

    if (loading) {
        return (
            <div className="bg-white shadow-sm rounded-lg border">
                <div className="px-6 py-4 border-b border-gray-200">
                    <h3 className="text-lg font-medium text-gray-900">
                        Liste des stages
                    </h3>
                </div>
                <div className="p-8 text-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-4"></div>
                    <p className="text-gray-600">Chargement des stages...</p>
                </div>
            </div>
        );
    }

    if (stages.length === 0) {
        return (
            <div className="bg-white shadow-sm rounded-lg border">
                <div className="px-6 py-4 border-b border-gray-200">
                    <h3 className="text-lg font-medium text-gray-900">
                        Liste des stages (0)
                    </h3>
                </div>
                <div className="p-8 text-center">
                    <p className="text-gray-500">Aucun stage trouvé</p>
                </div>
            </div>
        );
    }

    return (
        <div className="bg-white shadow-sm rounded-lg border">
            <div className="px-6 py-4 border-b border-gray-200">
                <h3 className="text-lg font-medium text-gray-900">
                    Liste des stages ({stages.length})
                </h3>
            </div>

            <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                    <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Entreprise
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Poste
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Commune
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Filière
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Parcours
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Niveau
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Famille Métiers
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Capacité
                        </th>
                        <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Actions
                        </th>
                    </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                    {stages.map((stage) => {
                        const available_places = Math.round(stage.stage_offer.capacity_total - stage.stage_offer.capacity_filled)
                        const capacityStyles = getCapacityStyles(
                            stage.stage_offer.capacity_filled,
                            stage.stage_offer.capacity_total,
                            available_places
                        );
                        const progressPercentage = stage.stage_offer.capacity_total > 0
                            ? (stage.stage_offer.capacity_filled / stage.stage_offer.capacity_total) * 100
                            : 0;

                        return (
                            <tr key={stage.id} className="hover:bg-gray-50">
                                <td className="px-4 py-4 whitespace-nowrap">
                                    <div className="text-sm font-medium text-gray-900">
                                        {stage.stage_offer.enterprise}
                                    </div>
                                    <div className="text-xs text-gray-500">
                                        {stage.stage_offer.sector || 'Secteur non renseigné'}
                                    </div>
                                </td>
                                <td className="px-4 py-4">
                                    <div className="text-sm text-gray-900 max-w-xs truncate" title={stage.poste}>
                                        {stage.stage_offer.position}
                                    </div>
                                </td>
                                <td className="px-4 py-4 whitespace-nowrap">
                                    <div className="text-sm text-gray-900">
                                        {stage.stage_offer.address || 'Non renseigné'}
                                    </div>
                                </td>
                                <td className="px-4 py-4 whitespace-nowrap">
                                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getFiliereStyles(stage.filiere)}`}>
                                            {stage.stage_offer.filiere}
                                        </span>
                                </td>
                                <td className="px-4 py-4 whitespace-nowrap">
                                        <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full border ${getParcoursStyles(stage.parcours || '')}`}>
                                            {stage.stage_offer.parcours || 'Non défini'}
                                        </span>
                                </td>
                                <td className="px-4 py-4 whitespace-nowrap">
                                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded ${getNiveauStyles(stage.niveau_scolaire || '')}`}>
                                            {stage.stage_offer.scolar_level || 'Non défini'}
                                        </span>
                                </td>
                                <td className="px-4 py-4">
                                    <div className="text-sm text-gray-900 max-w-xs truncate" title={stage.famille_metiers}>
                                        {stage.stage_offer.job_family || 'Non renseigné'}
                                    </div>
                                </td>
                                <td className="px-4 py-4 whitespace-nowrap">
                                    <div className="text-sm text-gray-900">
                                            <span className={`font-medium ${capacityStyles.text}`}>
                                                {stage.stage_offer.capacity_filled}/{stage.stage_offer.capacity_total}
                                            </span>
                                        <span className="text-gray-500 ml-1 text-xs block">
                                                ({available_places || 0} libres)
                                            </span>
                                    </div>
                                    <div className="mt-1 w-full bg-gray-200 rounded-full h-1.5">
                                        <div
                                            className={`h-1.5 rounded-full ${capacityStyles.bar}`}
                                            style={{ width: `${progressPercentage}%` }}
                                        />
                                    </div>
                                </td>
                                <td className="px-4 py-4 whitespace-nowrap text-right text-sm font-medium">
                                    <div className="flex justify-end">
                                        <button
                                            onClick={() => onEdit(stage)}
                                            className="inline-flex items-center px-3 py-1 bg-blue-100 text-blue-700 rounded-md hover:bg-blue-200 transition-colors text-sm"
                                            title="Modifier le stage"
                                        >
                                            <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                            </svg>
                                            Modifier
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        );
                    })}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default StagesList;
