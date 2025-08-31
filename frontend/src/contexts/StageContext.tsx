'use client'

import React, { createContext, useContext, useReducer, useEffect, ReactNode } from 'react'
import { Stage, StagesData } from '@/types/stage'
import { getStages, addStage } from '@/lib/stageApi'

interface StageState {
    stages: Stage[]
    loading: boolean
    error: string | null
    filters: {
        commune?: string
        filiere?: string
        sector?: string
    }
}

type StageAction =
    | { type: 'FETCH_START' }
    | { type: 'FETCH_SUCCESS'; payload: Stage[] }
    | { type: 'FETCH_ERROR'; payload: string }
    | { type: 'ADD_STAGE'; payload: Stage }
    | { type: 'SET_FILTERS'; payload: Partial<StageState['filters']> }

interface StageContextType extends StageState {
    refetchStages: () => Promise<void>
    addNewStage: (stage: Omit<Stage, 'id'>) => Promise<void>
    setFilters: (filters: Partial<StageState['filters']>) => void
    filteredStages: Stage[]
}

const StageContext = createContext<StageContextType | undefined>(undefined)

function stageReducer(state: StageState, action: StageAction): StageState {
    switch (action.type) {
        case 'FETCH_START':
            return { ...state, loading: true, error: null }

        case 'FETCH_SUCCESS':
            return { ...state, loading: false, stages: action.payload, error: null }

        case 'FETCH_ERROR':
            return { ...state, loading: false, error: action.payload }

        case 'ADD_STAGE':
            return { ...state, stages: [...state.stages, action.payload] }

        case 'SET_FILTERS':
            return { ...state, filters: { ...state.filters, ...action.payload } }

        default:
            return state
    }
}

export function StageProvider({ children }: { children: ReactNode }) {
    const [state, dispatch] = useReducer(stageReducer, {
        stages: [],
        loading: true,
        error: null,
        filters: {}
    })

    // Fetch initial des données
    const refetchStages = async () => {
        dispatch({ type: 'FETCH_START' })
        try {
            const data: StagesData = await getStages()
            dispatch({ type: 'FETCH_SUCCESS', payload: data.stages })
        } catch (error) {
            dispatch({ type: 'FETCH_ERROR', payload: error instanceof Error ? error.message : 'Erreur inconnue' })
        }
    }

    // Ajouter un nouveau stage
    const addNewStage = async (stageData: Omit<Stage, 'id'>) => {
        try {
            const newStage = await addStage(stageData)
            dispatch({ type: 'ADD_STAGE', payload: newStage })
        } catch (error) {
            throw error // Re-throw pour que le composant puisse gérer l'erreur
        }
    }

    // Mettre à jour les filtres
    const setFilters = (filters: Partial<StageState['filters']>) => {
        dispatch({ type: 'SET_FILTERS', payload: filters })
    }

    // Stages filtrés (calculés à chaque render)
    const filteredStages = state.stages.filter(stage => {
        if (state.filters.commune && stage.commune !== state.filters.commune) return false
        if (state.filters.filiere && stage.filiere !== state.filters.filiere) return false
        if (state.filters.sector && stage.sector !== state.filters.sector) return false
        return true
    })

    // Charge les données au montage
    useEffect(() => {
        refetchStages()
    }, [])

    return (
        <StageContext.Provider value={{
            ...state,
            refetchStages,
            addNewStage,
            setFilters,
            filteredStages
        }}>
            {children}
        </StageContext.Provider>
    )
}

export const useStages = () => {
    const context = useContext(StageContext)
    if (!context) {
        throw new Error('useStages must be used within a StageProvider')
    }
    return context
}
