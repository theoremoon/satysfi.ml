module Main exposing (main)

import Browser
import Browser.Navigation
import Html exposing (..)
import Html.Attributes exposing (attribute, class, src, type_, value)
import Html.Events exposing (onClick, onInput)
import Http
import Url exposing (Url)



-- MAIN


main =
    Browser.application
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        , onUrlRequest = OnUrlRequest
        , onUrlChange = OnUrlChange
        }



-- MODEL


type alias Model =
    { currentSource : String
    , product : Maybe String
    }


type Msg
    = OnUrlChange Url
    | OnUrlRequest Browser.UrlRequest
    | CompileRequest
    | CompileResult (Result Http.Error String)
    | SourceUpdate String



-- INIT


init : () -> Url -> Browser.Navigation.Key -> ( Model, Cmd Msg )
init _ url navKey =
    ( Model "" Nothing
    , Cmd.none
    )



-- UPDATE


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SourceUpdate newSource ->
            ( { model | currentSource = newSource }, Cmd.none )

        CompileRequest ->
            ( model, compileRequest model.currentSource )

        CompileResult (Result.Ok result) ->
            ( { model | product = Just result }, Cmd.none )

        CompileResult (Result.Err _) ->
            ( model, Cmd.none )

        _ ->
            ( model, Cmd.none )



-- CMDS


compileRequest : String -> Cmd Msg
compileRequest source =
    Http.post
        { url = "/compile"
        , body = Http.stringBody "text/plain" source
        , expect = Http.expectString CompileResult
        }



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Browser.Document Msg
view model =
    { title = "Hello World"
    , body =
        [ div
            [ class "container"
            ]
            [ div [ class "menu" ]
                [ button [ onClick CompileRequest ] [ text "COMPILE" ]
                ]
            , div [ class "main" ]
                [ satysfiEditor [] model
                , productViewer [] model
                ]
            ]
        ]
    }


satysfiEditor : List (Attribute Msg) -> Model -> Html Msg
satysfiEditor attrs model =
    textarea
        (attrs
            ++ [ class "editor"
               , value model.currentSource
               , onInput SourceUpdate
               ]
        )
        []


productViewer : List (Attribute Msg) -> Model -> Html Msg
productViewer attrs model =
    case model.product of
        Just pdf ->
            embed
                (attrs
                    ++ [ src ("data:application/pdf;base64," ++ pdf) ]
                )
                []

        Nothing ->
            div attrs [ text "Nothing there" ]
