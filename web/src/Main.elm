module Main exposing (main)

import Base64
import Browser
import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (attribute, class, href, src, type_, value)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode as D
import Json.Encode as E
import Url exposing (Url)
import Url.Builder
import Url.Parser as Parser exposing ((</>), Parser)



-- MAIN


main =
    Browser.application
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        , onUrlRequest = LinkClicked
        , onUrlChange = UrlChanged
        }



-- ROUTE


type Route
    = IndexPage
    | RegisterPage


routeParser : Parser (Route -> a) a
routeParser =
    Parser.oneOf
        [ Parser.map IndexPage Parser.top
        , Parser.map RegisterPage (Parser.s "register")
        ]



-- MODEL


type ViewModel
    = Index
    | Register
    | Notfound


type alias Model =
    { key : Nav.Key
    , url : Url
    , model : ViewModel
    }


urlToModel : Url -> ViewModel
urlToModel url =
    case Parser.parse routeParser url of
        Just IndexPage ->
            Index

        Just RegisterPage ->
            Register

        Nothing ->
            Notfound



-- INIT


init : () -> Url -> Nav.Key -> ( Model, Cmd Msg )
init _ url key =
    ( Model key url (urlToModel url), Cmd.none )



-- UPDATE


type Msg
    = LinkClicked Browser.UrlRequest
    | UrlChanged Url


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        LinkClicked urlRequest ->
            case urlRequest of
                Browser.Internal url ->
                    ( model, Nav.pushUrl model.key (Url.toString url) )

                Browser.External href ->
                    ( model, Nav.load href )

        UrlChanged url ->
            ( { model | url = url, model = urlToModel url }, Cmd.none )



-- SUBSCRIPTION


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none



-- VIEW


view : Model -> Browser.Document Msg
view model =
    { title = "SATySFi Online"
    , body =
        [ div []
            [ h1 [] [ text "SATySFi Online" ]
            , div []
                [ text
                    (case model.model of
                        Index ->
                            "Index"

                        Register ->
                            "Register"

                        Notfound ->
                            "Notfound"
                    )
                ]
            , a [ href "/register" ] [ text "register" ]
            , a [ href "/takoyaki" ] [ text "takoyaki" ]
            ]
        ]
    }
